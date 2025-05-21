import { getPostgresConnection } from './db.js';
const db = await getPostgresConnection()

// console.log(`process ${process.pid} started.`);

process.on('message', (items) => {
  db.students.insert(items)
    .then(() => {
      process.send('item-done');
    })
    .catch((error) => {
      console.error(error);
    });
});
