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
    except Exception:
        import requests
        r = requests.get("https://restcountries.com/v3.1/all")
        data = r.json()
        with open("data/countries.json", "w") as f:
            json.dump(data, f)
        print("Local data not found, requesting from API...")
        return data


initial_data = get_data()


flag_emojis = []
for country in initial_data:
    flag_emojis.append(country["flag"])


def get_random_flag_emoji(flag_emojis):
    return random.choice(flag_emojis)


@app.route("/")
def index():
    return render_template("home.html")


@app.route("/search_continents", methods=["GET", "POST"])
def search_continents():
    if request.method == "GET":
        countries = initial_data
        continents = []
        for country in countries:
            for continent in country["continents"]:
                if continent not in continents:
                    continents.append(continent)

        emoji = get_random_flag_emoji(flag_emojis)
        return render_template(
            "search_continents.html", continents=continents, flag_emoji=emoji
        )
    if request.method == "POST":
        data = initial_data

        if request.form["continent"].lower() == "all" or \
           request.form["continent"].lower() == "":
            return data

        matched_data = []
        for country in data:
            if request.form["continent"].lower() in \
               [continent.lower() for continent in country["continents"]]:
                matched_data.append(country)
        return matched_data


@app.route("/country_guesser")
def guesser():
    data = initial_data
    countries = []
    for country in data:
        countries.append(country["name"]["common"])
    if request.method == "GET":
        emoji = get_random_flag_emoji(flag_emojis)
        return render_template(
            "country_guesser.html", countries=countries, flag_emoji=emoji
        )


@app.route("/capital_guesser")
def capital_guesser():
    data = initial_data
    capitals = []
    for country in data:
        try:
            for capital in country["capital"]:
                capitals.append(capital)
        except KeyError:
            capitals.append("None")

    emoji = get_random_flag_emoji(flag_emojis)
    return render_template(
        "capital_guesser.html", capitals=capitals, flag_emoji=emoji
    )


@app.route("/country")
def get_country():
    data = initial_data
    return random.choice(data)


if __name__ == "__main__":
    app.run(port=8000, debug=True, host=HOST)
