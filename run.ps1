$env:ORCH_URL = "http://orch:8080/internal/task"
$env:COMPUTING_POWER = 3
$env:TIME_ADDITION_MS = 10
$env:TIME_SUBTRACTION_MS = 10
$env:TIME_MULTIPLICATION_MS = 10
$env:TIME_DIVISION_MS = 10

docker-compose up --build --remove-orphans
