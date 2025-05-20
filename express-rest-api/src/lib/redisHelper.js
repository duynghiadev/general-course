import redisClient from "../config/redis.js";
import logger from "../utils/logger.js";

const getCache = async (key) => {
  try {
    const value = await redisClient.get(key);
    return value ? JSON.parse(value) : null;
  } catch (error) {
    logger.error("Error getting cache:", error);
    return null;
  }
};

const setCache = async (key, value, expirationInSeconds) => {
  try {
    const stringValue = JSON.stringify(value);
    await redisClient.set(key, stringValue, "EX", expirationInSeconds);
  } catch (error) {
    logger.error("Error setting cache:", error);
  }
};

const deleteCache = async (key) => {
  try {
    await redisClient.del(key);
  } catch (error) {
    logger.error("Error deleting cache:", error);
  }
};

// Get a specific field from a hash
const hGet = async (hash, field) => {
  try {
    const value = await redisClient.hget(hash, field);
    return value ? JSON.parse(value) : null;
  } catch (error) {
    logger.error("Error getting cache:", error);
    return null;
  }
};

// Set a field in a hash
const hSet = async (hash, field, value) => {
  try {
    await redisClient.hset(hash, field, JSON.stringify(value));
  } catch (error) {
    logger.error("Error setting cache:", error);
  }
};

// Delete a field from a hash
const hDel = async (hash, field) => {
  try {
    await redisClient.hdel(hash, field);
  } catch (error) {
    logger.error("Error deleting cache:", error);
  }
};

// Get all fields in a hash
const hGetAll = async (hash) => {
  try {
    const data = await redisClient.hgetall(hash);
    return Object.keys(data).length
      ? Object.fromEntries(
          Object.entries(data).map(([key, value]) => [key, JSON.parse(value)])
        )
      : null;
  } catch (error) {
    logger.error("Error getting all cache:", error);
    return null;
  }
};

export const expire = async (key, seconds) => {
  try {
    await redisClient.expire(key, seconds);
  } catch (error) {
    logger.error(`Error setting expiration for key ${key}:`, error);
  }
};

export {
  getCache,
  setCache,
  deleteCache,
  hGet,
  hSet,
  hDel,
  hGetAll,
  redisClient,
};
