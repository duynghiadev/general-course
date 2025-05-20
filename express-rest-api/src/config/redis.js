import Redis from "ioredis";
import logger from "../utils/logger.js";

const redisClient = new Redis({
  host: process.env.REDIS_HOST,
  port: process.env.REDIS_PORT,
});

redisClient.on("connect", () => {
  logger.info("✅ Successfully connected to Redis server");
});

redisClient.on("error", (err) => logger.error("Redis Client Error:", err));

export default redisClient;
