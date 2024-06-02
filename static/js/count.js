const form = document.getElementById("form-post-visit");
const count = document.getElementById("count");
async function sendData() {
    const formData = new FormData(form);

    try {
        const response = await fetch("/count", {
            method: "POST",
            body: formData,
            headers: {
                "Content-Type": "application/x-www-form-urlencoded"
            }
        });
        r = await response.text();
        count.innerText = r;
    } catch (e) {
        console.log(e);
    }
}

form.addEventListener("submit", (e) => {
    e.preventDefault();
    sendData();
});
