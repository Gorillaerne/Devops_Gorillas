import { callSearchRestApi } from "./api_calls.js"
import { checkIfLoggedIn } from "./reuseable_functions.js"

const searchButton = document.getElementById("search-button")
const searchInput = document.getElementById("search-input")
const resultsDiv = document.getElementById("results")

checkIfLoggedIn() 


function doSearch() {
    const query = searchInput.value
    callSearchRestApi(query)
        .then(searchResults => {
            renderSearchResults(searchResults.data)
        })
        .catch(error => {
            console.error(error)
        })
}

searchButton.addEventListener("click", doSearch)

searchInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") doSearch()
})

function renderSearchResults(searchResults) {

    resultsDiv.innerHTML = ""

    if (!searchResults || searchResults.length === 0) {
        const noResults = document.createElement("p")
        noResults.classList.add("no-results")
        noResults.textContent = "No result on search"
        resultsDiv.appendChild(noResults)
        return
    }

    searchResults.forEach(searchResult => {
        const item = document.createElement("div")
        item.classList.add("result-item")

        const link = document.createElement("a")
        link.href = searchResult.URL

        const titleSpan = document.createElement("span")
        titleSpan.classList.add("result-title")
        titleSpan.textContent = searchResult.title

        const snippet = document.createElement("p")
        snippet.classList.add("result-snippet")
        snippet.textContent = searchResult.description || "Ingen beskrivelse tilgængelig"

        link.appendChild(titleSpan)
        item.appendChild(link)
        item.appendChild(snippet)
        resultsDiv.appendChild(item)

    })
}





