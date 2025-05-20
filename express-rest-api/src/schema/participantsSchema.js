import mongoose from "mongoose";
import { PARTICIPANTS_ROLE } from "../config/constants.js";

const participantsSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
  },
  age: {
    type: Number,
    required: true,
  },
  role: {
    type: String,
    required: true,
    enum: PARTICIPANTS_ROLE,
  },
});

const Participants = mongoose.model("Participants", participantsSchema);
export default Participants;
