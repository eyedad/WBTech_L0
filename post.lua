wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"


local file = io.open("test_files/8.json", "r")
wrk.body = file:read("*a")
file:close()
