# Dissertation Project




The bin/run.sh file runs the experiment by calling a docker-compose tool that composes the target system.
After all the Docker images start, UI will be accessible through on 8081 port. In the case of the local execution, it will be http://127.0.0.1:8081/.
The shutdown is made by calling bin/shutdown.sh.

Example:
```
cd bin
./run.sh
```
or
```
cd bin
./shutdown.sh
```