
const imageForm = document.querySelector("#imageForm")
const imageInput = document.querySelector("#imageInput")

imageForm.addEventListener("submit", async e => {
    e.preventDefault()
    const file = imageInput.files[0]

    // get secure url from our server
    const signed = await fetch("http://127.0.0.1:8022/s3Url").then(res => res.json())
    const url = signed.url

    console.log("s3 url", url)

    // post image direct to s3
    await fetch(url, {
        method: "PUT",
        body: file
    })
    // post request to my server with url
    const imageUrl = url.split("?")[0]
    console.log("imageUrl", imageUrl)
})

