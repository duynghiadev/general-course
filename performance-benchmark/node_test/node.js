import express from "express";
const app = express();

app.get("/ping", (req, res) => res.send("pong"));

app.listen(3000, () => console.log("Node.js listening on 3000"));
