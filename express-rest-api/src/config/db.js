import mongoose from "mongoose";
import logger from "../utils/logger.js";

const connect = async (dbname) => {
  try {
    let connectString = process.env.MONGO_URI || "";
    if (connectString === "") {
      throw new Error("Please provide a valid connection string");
    }

    connectString = connectString.replace("<dbname>", dbname);
    await mongoose.connect(connectString);

    logger.info(`✅ Successfully connected to the ${dbname} database`);
  } catch (err) {
    logger.error(
      `❌ Couldn't connect to the ${dbname} database: ${err?.message}`
    );
    process.exit(1);
  }
};

export default connect;
