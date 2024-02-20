# Backup Policy
Every expression is saved completely in MongoDB in 2 cases:
 * Every 10 seconds in order to avoid losing data (autosave)
 * When restarting the server manually

# Server status codes:
 * 200, everything is fine
 * 400, expression isn't valid (division by zero or unsolvable expressions)
 * 500, smth is wrong on backend
 * 502, too much requests (there are 3 independent servers, and each one could handle only 10 expressions at the same time, so **30** expressions is the limit by default, you can change settings in `.env`)
 * 503, expression "expired"

---

# Turning off the server
Turning off the server can be done by visiting `localhost:8000/kill`. 

---

# Checking if expression is valid or not
The expression should meet the following requirements:
 * Expression contains only numbers (0-9), operators (+, -, *, /) and round brackets "()"
 * Expression contains right amount of round brackets (e.g. `2+(2` is invalid)
 * Expression is solvable (no division by 0)
 * There are no negative or float numbers
They are checked by various ways:
 * 1st one is checked using regexp both on frontend and backend
 * The 2nd one is checked in the proccess of parsing
 * The 3rd one is checked while solving the expression on the backend
 * The 4th one is also checked while solving the expression on the backend. Expression is invalid if there are negative numbers (counting with negative numbers makes server superclusters go boom-boom in alternative reality). All the float numbers are getting transformed to integer parts just by cutting their fractional part off (flooring). 

---

# Launching the server
Just run the `main.go` file! (also edit .env to add your MongoDB URI in it and run `go mod tidy` to install dependencies)
**WARNING: Having own MongoDB cluster is essential for launching**

---

# Roadmap 
* 19.02 - 20.02.24: making frontend
* 21.02 - 22.02.24: adding an option to resolve expression on different servers at the single moment