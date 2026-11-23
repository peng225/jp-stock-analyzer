import { defineConfig, loadEnv } from 'vite';
import fs from 'node:fs';
import path from 'node:path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  // .env ファイルを作るなどして設定が必要。
  const dataDir = env.STOCK_DATA_DIR;

  return {
    plugins: [
      {
        name: 'stock-data-api',

        configureServer(server) {
          server.middlewares.use('/api/csv', (req, res) => {
            if (req.method !== 'GET') {
              res.statusCode = 405;
              res.end('Method Not Allowed');
              return;
            }

            if (!dataDir) {
              res.statusCode = 500;
              res.end('STOCK_DATA_DIR is not configured');
              return;
            }

            try {
              const files = fs.readdirSync(dataDir)
                .filter(file => file.toLowerCase().endsWith('.csv'))
                .sort();

              const result = files.map(file => ({
                name: file,
                content: fs.readFileSync(
                  path.join(dataDir, file),
                  'utf-8'
                ),
              }));

              res.setHeader('Content-Type', 'application/json; charset=utf-8');
              res.end(JSON.stringify(result));

            } catch (error) {
              console.error(error);

              res.statusCode = 500;
              res.end('Failed to read CSV files');
            }
          });
        },
      },
    ],
  };
});
