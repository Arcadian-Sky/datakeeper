db.auth(
    process.env.MONGO_INITDB_ROOT_USERNAME,
    process.env.MONGO_INITDB_ROOT_PASSWORD
)

db = db.getSiblingDB(process.env.MONGO_APP_DB)

db.createUser({
    user: process.env.MONGO_APP_USERNAME,
    pwd: process.env.MONGO_APP_PASSWORD,
    roles: [
        {
            role: 'root',
            db: process.env.MONGO_APP_DB,
        },
    ],
});