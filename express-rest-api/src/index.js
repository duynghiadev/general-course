import express from "express";
import dotenv from "dotenv";
import logger from "./utils/logger.js";

// Controllers
import participantRouter from "./controller/participantsController.js";

// Database connection
import { DB_NAME } from "./config/constants.js";
import connect from "./config/db.js";
dotenv.config();
await connect(DB_NAME);

const app = express();
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

const port = process.env.PORT || 8000;

app.listen(port, () => {
  logger.info(`Server is running on port ${port}`);
  logger.info(
    `Selected environment is ${process.env.NODE_ENV || "development"}`
  );
});

// Routes
app.use("/participants", participantRouter);

app.get("/ping", (req, res) => {
  res.send("pong");
});

// Jobs
import "./job/cacheRefresherJob.js";
