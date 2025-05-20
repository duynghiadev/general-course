import express from "express";
import dotenv from "dotenv";
import { DB_NAME } from "./config/constants.js";
import connect from "./config/db.js";

dotenv.config();

await connect(DB_NAME);

const app = express();
const port = process.env.PORT || 8000;

app.listen(port, () => console.log(`Server is running on port ${port}`));

app.get("/ping", (req, res) => {
  res.send("pong");
});

app.get("/", (req, res) => {
  res.send("Hello World!");
});
