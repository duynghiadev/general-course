import mongoose from "mongoose";

const connect = async (dbname) => {
  try {
    let connectString = process.env.MONGO_URI || "";
    if (connectString === "") {
      throw new Error("Please provide a valid connection string");
    }

    connectString = connectString.replace("<dbname>", dbname);
    await mongoose.connect(connectString);
    console.log(`Successfully connected to the ${dbname} database`);
  } catch (err) {
    console.log(`Couldn't connect to the ${dbname} database`, err?.message);
    process.exit(1);
  }
};

export default connect;
