const form = document.getElementById("fetch-continents-form");
const formSort = document.getElementById("sort-method");
const countryContainer = document.querySelector("#country-container");

function generateFlagHTML(country) {
    return `
        <div class="country-item">
            <i class="country-flag-icon">
                ${country.flag}
            </i>
            <p class="country-name">
                ${country.name.common}
            </p>
            <p class="country-population">
                Population: ${country.population.toLocaleString()}
            </p>
            <p class="country-flag-description">
                ${country.flags.alt ? country.flags.alt : "No Description Available"}
            </p>
        </div>
    `
}

async function displayCountryData(data) {
    const countryItems = data.reduce((acc, curr) => {
        return acc += generateFlagHTML(curr)
    },'')

    countryContainer.innerHTML = countryItems
}

async function fetchCountryData(formData, countryData) {
    try {
        const response = await fetch("http://localhost:8000/continents", {
            method: "POST",
            body: formData,
        });
        countryData = await response.json();
        return countryData;
    } catch(e) {
        console.error(e);
    }
}

function sortContinentData(formData, countryData) {
    if (formData.get("sort-method") === "alpha") {
        countryData.sort((a, b) => {
            if (a.name.common > b.name.common) {
                return 1;
            }
            if (a.name.common < b.name.common) {
                return -1;
            }
            return 0;
        })
    } else if (formData.get("sort-method") === "pop") {
        countryData.sort((a, b) => {
            return b.population - a.population;
        })
    }
    return countryData;
}

form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const formData = new FormData(form);
    let countryData = await fetchCountryData(formData);

    countryData = sortContinentData(formData, countryData);

    displayCountryData(countryData);
});

formSort.addEventListener("change", async (e) => {
    e.preventDefault();
    const formData = new FormData(form);
    let countryData = await fetchCountryData(formData);

    countryData = sortContinentData(formData, countryData);

    displayCountryData(countryData);
});
