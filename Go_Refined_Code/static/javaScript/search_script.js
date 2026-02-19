import { callSearchRestApi } from "./api_calls.js"

const searchButton = document.getElementById("search-button")
const searchInput = document.getElementById("search-input")
const resultsDiv = document.getElementById("results")




searchButton.addEventListener("click", function () {

    const query = searchInput.value

    callSearchRestApi(query)
        .then(searchResults => {
            renderSearchResults(searchResults.data)
        })
        .catch(error => {
            console.error(error)
        })
})



function renderSearchResults(searchResults) {

    resultsDiv.innerHTML = ""
    console.log(searchResults)
    searchResults.forEach(searchResult => {
        console.log(searchResult)
        const searchResultDiv = document.createElement("div")
        const searchResultH2 = document.createElement("H2")
        const searchResultA = document.createElement("a")
        searchResultA.textContent = searchResult.title
        searchResultA.classList.add("search-result-title")
        searchResultA.href = searchResult.URL
        searchResultH2.appendChild(searchResultA)

        searchResultDiv.appendChild(searchResultH2)

        resultsDiv.appendChild(searchResultDiv)


    })

}





