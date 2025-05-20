// TODO: Lift this to env variables
import {
  hGet,
  hSet,
  hDel,
  hGetAll,
  expire,
  redisClient,
} from "../lib/redisHelper.js";
import logger from "../utils/logger.js";

const PARTICIPANTS_KEY = "participants"; // HashMap key
const PARTICIPANTS_EXPIRATION = 60 * 60; // 1 hour

const getParticipantCache = async (id) => {
  return await hGet(PARTICIPANTS_KEY, id);
};

const setParticipantCache = async (id, value) => {
  await hSet(PARTICIPANTS_KEY, id, value);
};

const delParticipantCache = async (id) => {
  await hDel(PARTICIPANTS_KEY, id);
};

const getAllParticipantsCache = async () => {
  return await hGetAll(PARTICIPANTS_KEY);
};

const setAllParticipantsCache = async (participants) => {
  try {
    if (!Array.isArray(participants) || participants.length === 0) {
      logger.warn("⚠️  No participants to cache.");
      return;
    }

    const pipeline = redisClient.pipeline();

    participants.forEach((participant) => {
      const data = participant.toObject ? participant.toObject() : participant; // Convert Mongoose doc to plain object
      pipeline.hset(
        PARTICIPANTS_KEY,
        data._id.toString(),
        JSON.stringify(data)
      );
    });

    pipeline.expire(PARTICIPANTS_KEY, PARTICIPANTS_EXPIRATION);

    await pipeline.exec();

    const expirationTime = new Date(
      Date.now() + PARTICIPANTS_EXPIRATION * 1000
    ).toISOString();
    logger.info(`✅ Participants cache updated. Expires at ${expirationTime}.`);
  } catch (error) {
    logger.error("❌ Error setting participants cache:", error);
  }
};

export {
  getParticipantCache,
  setParticipantCache,
  delParticipantCache,
  getAllParticipantsCache,
  setAllParticipantsCache,
};
