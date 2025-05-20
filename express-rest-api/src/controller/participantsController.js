import { Router } from "express";
import * as participantsService from "../service/participantsService.js";

const participantRouter = Router();

participantRouter.get("/", participantsService.getAllParticipants);
participantRouter.get("/:id", participantsService.getParticipant);
participantRouter.post("/", participantsService.createParticipant);
participantRouter.put("/:id", participantsService.updateParticipant);
participantRouter.delete("/:id", participantsService.deleteParticipant);

export default participantRouter;
