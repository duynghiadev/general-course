import {
  getAll,
  getById,
  create,
  update,
  remove,
} from "../repository/participantsRepository.js";

import {
  getParticipantCache,
  setParticipantCache,
  delParticipantCache,
  getAllParticipantsCache,
} from "../cacheHelper/participantCacheHelper.js";

export const getAllParticipants = async (req, res) => {
  try {
    // Check if participants exist in cache
    const cachedParticipants = await getAllParticipantsCache();
    if (cachedParticipants) {
      return res.json(Object.values(cachedParticipants)); // Convert HashMap values to array
    }

    const participants = await getAll();

    // Cache update (store each participant in HashMap)
    for (const participant of participants) {
      await setParticipantCache(participant._id, participant);
    }

    res.json(participants);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

export const getParticipant = async (req, res) => {
  const { id } = req.params;
  try {
    // Check if participant exists in cache
    const cachedParticipant = await getParticipantCache(id);
    if (cachedParticipant) {
      return res.json(cachedParticipant);
    }

    const participant = await getById(id);
    if (!participant) {
      return res.status(404).json({ error: "Participant not found" });
    }

    // Cache update
    await setParticipantCache(participant._id, participant);

    res.json(participant);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

export const createParticipant = async (req, res) => {
  const { name, age, role } = req.body;
  try {
    const participant = await create({ name, age, role });

    // Cache update
    await setParticipantCache(participant._id, participant);

    res.status(201).json(participant);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

export const updateParticipant = async (req, res) => {
  const { id } = req.params;
  const { name, age, role } = req.body;
  try {
    const participant = await update(id, { name, age, role });
    if (!participant) {
      return res.status(404).json({ error: "Participant not found" });
    }

    // Cache update
    await setParticipantCache(participant._id, participant);

    res.json(participant);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};

export const deleteParticipant = async (req, res) => {
  const { id } = req.params;
  try {
    const participant = await remove(id);
    if (!participant) {
      return res.status(404).json({ error: "Participant not found" });
    }

    // Cache update
    await delParticipantCache(id);

    res
      .status(200)
      .json({ message: "Participant deleted successfully", participant });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
};
