import Queue from "bull";
import logger from "../utils/logger.js";
import { getAll } from "../repository/participantsRepository.js";
import { setAllParticipantsCache } from "../cacheHelper/participantCacheHelper.js";

// Create a Bull queue
const participantQueue = new Queue("participant-cache-refresh", {
  redis: {
    host: process.env.REDIS_HOST,
    port: process.env.REDIS_PORT,
  },
});

// Job processor: This function will run each time a job is processed
participantQueue.process(async (job) => {
  try {
    logger.info("🔄 Refreshing participants cache...");
    const participants = await getAll();

    await setAllParticipantsCache(participants);

    logger.info("✅ Participants cache updated successfully");
  } catch (error) {
    logger.error(`❌ Error refreshing participants cache: ${error.message}`);
  }
});

const scheduleCacheRefreshJob = () => {
  logger.info("⏳ Scheduling cache refresher job...");
  participantQueue.add({}, { repeat: { cron: "0 * * * *" } }); // Runs every hour
  logger.info(`🔔 Cache refresher job scheduled.`);
};

export { participantQueue, scheduleCacheRefreshJob };
