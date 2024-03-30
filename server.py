from flask import Flask
from flask import render_template
from flask import request
import json
import random

HOST = "localhost"

app = Flask(__name__)
app.config["MIME_TYPES"] = {
    ".js": "application/javascript",
}


def get_data():
    try:
        with open("data/countries.json", "r") as f:
            data = json.load(f)
        print("Loading local data...")
        return data
    except:
        import requests
        r = requests.get("https://restcountries.com/v3.1/all")
        data = r.json()
        with open("data/countries.json", "w") as f:
            json.dump(data, f)
        print("Local data not found, requesting from API...")
        return data


@app.route("/")
def index():
    countries = get_data()
    continents = []
    for country in countries:
        for continent in country["continents"]:
            if continent not in continents:
                continents.append(continent)

    return render_template("index.html", continents=continents)


@app.route("/continents", methods=["POST"])
def continents():
    if request.method == "POST":
        data = get_data()

        if request.form["continent"].lower() == "all":
            return data

        matched_data = []
        for country in data:
            if request.form["continent"] in country["continents"]:
                matched_data.append(country)
        return matched_data


@app.route("/guesser")
def guesser():
    data = get_data()
    countries = []
    for country in data:
        countries.append(country["name"]["common"])
    if request.method == "GET":
        return render_template("guesser.html", countries=countries)


@app.route("/country")
def get_country():
    data = get_data()
    return random.choice(data)


if __name__ == "__main__":
    app.run(port=8000, debug=True, host=HOST)
