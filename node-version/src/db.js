import { MongoClient } from 'mongodb';
import pg from 'pg';
import format from 'pg-format';
const { Client } = pg;
// Connection URL for MongoDB

// Connection pool for PostgreSQL

async function getMongoConnection() {
    const mongoUrl = 'mongodb://root:example@localhost:27017';
    const client = new MongoClient(mongoUrl);

    try {
        await client.connect();

        const dbName = 'school';
        const db = client.db(dbName);
        const collection = db.collection('students');

        return {
            students: collection,
            client,
        };
    } catch (error) {
        console.error('Error connecting to MongoDB:', error);
        throw error;
    }
}

async function getPostgresConnection() {

    const client = new Client({
        user: 'erickwendel',
        host: 'localhost',
        database: 'school',
        password: 'mypassword',
        port: 5432,
    });

    await client.connect();
    return {
        client,
        students: {
            async insert(persons) {
              const values = persons.map(person => {
                const { name, email, age, registeredAt } = person;

                return [name, email, age, registeredAt];
              })
                const query = format('INSERT INTO students (name, email, age, registered_at) VALUES %L', values);

                await client.query(query);

            },
            async list(limit = 100) {
                const query = 'SELECT * FROM students LIMIT $1';
                const values = [limit];

                const result = await client.query(query, values);
                return result.rows;

            },
            async count() {
                const query = 'SELECT COUNT(*) as total FROM students';

                const result = await client.query(query);
                return Number(result.rows[0].total);

            },
            async deleteAll() {
                const query = 'DELETE FROM students';

                await client.query(query);
            },
            async createTable() {
                const createStudentsTableQuery = `
                        CREATE TABLE IF NOT EXISTS students (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            email VARCHAR(255) NOT NULL,
                            age INT NOT NULL,
                            registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                        )`;
                await client.query(createStudentsTableQuery);

            }
        }
    };
}

export { getMongoConnection, getPostgresConnection };
