#!/bin/bash
docker container rm minerva_tui -f
exec docker container run -it --name minerva_tui luksamuk/minerva_tui:latest

