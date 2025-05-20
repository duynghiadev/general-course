import mongoose from "mongoose";
import Participants from "../schema/participantsSchema.js";

export const getAll = async () => {
  return await Participants.find();
};

export const getById = async (id) => {
  if (!mongoose.Types.ObjectId.isValid(id)) {
    throw new Error("Invalid ID format");
  }
  return await Participants.findById(id);
};

export const create = async ({ name, age, role }) => {
  return await Participants.create({ name, age, role });
};

export const update = async (id, { name, age, role }) => {
  if (!mongoose.Types.ObjectId.isValid(id)) {
    throw new Error("Invalid ID format");
  }
  return await Participants.findByIdAndUpdate(
    id,
    { name, age, role },
    { new: true }
  );
};

export const remove = async (id) => {
  if (!mongoose.Types.ObjectId.isValid(id)) {
    throw new Error("Invalid ID format");
  }
  return await Participants.findByIdAndDelete(id);
};
