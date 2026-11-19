/// <reference types="vite/client" />
import Chart from 'chart.js/auto';
import Papa from 'papaparse';
import './style.css';

// CSVのデータ形式
interface StockData {
  code: string;
  year: string;

  sales: number | null;
  operatingProfit: number | null;
  ordinaryProfit: number | null;
  netProfit: number | null;
}

// グラフ表示項目
type ItemKey =
  | 'sales'
  | 'operatingProfit'
  | 'ordinaryProfit'
  | 'netProfit';

interface ItemDefinition {
  key: ItemKey;
  label: string;
  column: string;
}

// CSVの列名と内部データの対応
const items: ItemDefinition[] = [
  {
    key: 'sales',
    label: '売上高',
    column: '売上高',
  },
  {
    key: 'operatingProfit',
    label: '営業利益',
    column: '営業利益',
  },
  {
    key: 'ordinaryProfit',
    label: '経常利益',
    column: '経常利益',
  },
  {
    key: 'netProfit',
    label: '純利益',
    column: '純利益',
  },
];

// DOM
const fileInput =
  document.querySelector<HTMLInputElement>('#fileInput')!;

const stockSelect =
  document.querySelector<HTMLSelectElement>('#stockSelect')!;

const checkboxes =
  document.querySelectorAll<HTMLInputElement>(
    '.checkboxes input[type="checkbox"]'
  );

const statusElement =
  document.querySelector<HTMLDivElement>('#status')!;

const canvas =
  document.querySelector<HTMLCanvasElement>('#chart')!;

// 読み込んだデータ
let stockData: StockData[] = [];

// グラフインスタンス
let chart: Chart | undefined;

// CSVの数値を変換
function parseNumber(value: string | undefined): number | null {
  if (value === undefined || value.trim() === '' || value.trim() === '-') {
    return null;
  }

  const number = Number(value);

  return Number.isFinite(number) ? number : null;
}

// CSVの1行をStockDataに変換
function convertRow(row: string[]): StockData | null {
    if (row.length < 6) {
        return null;
    }

    const code = row[0]?.trim();
    const year = row[1]?.trim();

    if (!code || !year) {
        return null;
    }

    return {
        code,
        year,
        sales: parseNumber(row[2]),
        operatingProfit: parseNumber(row[3]),
        ordinaryProfit: parseNumber(row[4]),
        netProfit: parseNumber(row[5]),
    };
}

// 年月をソートするための値
function yearValue(year: string): number {
  const normalized = year.replace('/', '-');

  return Date.parse(`${normalized}-01`);
}

// 年月表示用
function formatYear(year: string): string {
  return year.replace('/', '年') + '月';
}

// CSVファイルを読み込む
async function loadFile(file: File): Promise<StockData[]> {
    return new Promise((resolve, reject) => {
        Papa.parse<string[]>(file, {
            header: false,
            skipEmptyLines: true,

            complete: (result) => {
                const rows = result.data
                    .map(convertRow)
                    .filter((row): row is StockData => row !== null);

                resolve(rows);
            },

            error: (error) => {
                reject(error);
            },
        });
    });
}

// 複数CSVファイルを読み込む
async function loadFiles(files: FileList): Promise<void> {
  statusElement.textContent = '読み込み中...';

  try {
    const allRows: StockData[] = [];

    for (const file of files) {
      const rows = await loadFile(file);
      allRows.push(...rows);
    }

    // 同じ銘柄・年度が複数ある場合は後から読み込んだデータを優先
    const dataMap = new Map<string, StockData>();

    for (const row of allRows) {
      const key = `${row.code}_${row.year}`;
      dataMap.set(key, row);
    }

    stockData = Array.from(dataMap.values());

    // コード順に並べる
    stockData.sort((a, b) => a.code.localeCompare(b.code));

    updateStockSelect();

    statusElement.textContent =
      `${files.length}ファイル、${stockData.length}件のデータを読み込みました。`;

    updateChart();
  } catch (error) {
    console.error(error);
    statusElement.textContent = 'CSVの読み込みに失敗しました。';
  }
}

// 銘柄選択プルダウンを更新
function updateStockSelect(): void {
  const codes = [...new Set(stockData.map(row => row.code))]
    .sort((a, b) => a.localeCompare(b));

  stockSelect.innerHTML = '';

  for (const code of codes) {
    const option = document.createElement('option');

    option.value = code;
    option.textContent = code;

    stockSelect.appendChild(option);
  }

  stockSelect.disabled = codes.length === 0;
}

// 選択されている表示項目
function getSelectedItems(): ItemDefinition[] {
  const selected = new Set(
    Array.from(checkboxes)
      .filter(checkbox => checkbox.checked)
      .map(checkbox => checkbox.value)
  );

  return items.filter(item => selected.has(item.label));
}

// グラフ更新
function updateChart(): void {
  const selectedCode = stockSelect.value;

  if (!selectedCode) {
    chart?.destroy();
    chart = undefined;
    return;
  }

  // 選択された銘柄のデータを取得
  const rows = stockData
    .filter(row => row.code === selectedCode)
    .sort((a, b) => yearValue(a.year) - yearValue(b.year));

  const selectedItems = getSelectedItems();

  // 横軸
  const labels = rows.map(row => formatYear(row.year));

  // グラフデータ
    const datasets = selectedItems.map(item => ({
        label: item.label,
        data: rows.map(row => row[item.key]),
        tension: 0.2,

        pointStyle:
            item.key === 'operatingProfit'
                ? 'circle'
                : item.key === 'ordinaryProfit'
                    ? 'rectRot'
                    : 'circle',

        pointRadius:
            item.key === 'ordinaryProfit'
                ? 7 : 5,
        pointHoverRadius:
            item.key === 'ordinaryProfit'
                ? 10 : 8,

        spanGaps: false,
    }));

  // 既存グラフを更新
  if (chart) {
    chart.data.labels = labels;
    chart.data.datasets = datasets;
    chart.update();
    return;
  }

  // 新規グラフを作成
  chart = new Chart(canvas, {
    type: 'line',

    data: {
      labels,
      datasets,
    },

    options: {
      responsive: true,
      maintainAspectRatio: false,

      interaction: {
        mode: 'index',
        intersect: false,
      },

      plugins: {
        legend: {
          position: 'top',
        },

        tooltip: {
          callbacks: {
            label: (context) => {
              const value = context.parsed.y;

              if (value === null) {
                return `${context.dataset.label}: データなし`;
              }

              return `${context.dataset.label}: ${
                value.toLocaleString('ja-JP')
              } 円`;
            },
          },
        },
      },

      scales: {
        x: {
          title: {
            display: true,
            text: '年度',
          },
        },

        y: {
          title: {
            display: true,
            text: '金額（円）',
          },

          ticks: {
            callback: (value) => {
              return Number(value).toLocaleString('ja-JP');
            },
          },
        },
      },
    },
  });
}

// イベント設定
fileInput.addEventListener('change', () => {
  if (fileInput.files && fileInput.files.length > 0) {
    loadFiles(fileInput.files);
  }
});

stockSelect.addEventListener('change', updateChart);

checkboxes.forEach(checkbox => {
  checkbox.addEventListener('change', updateChart);
});
