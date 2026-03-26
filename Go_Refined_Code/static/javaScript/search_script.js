import { callSearchRestApi } from "./api_calls.js"
import { checkIfLoggedIn } from "./reuseable_functions.js"

const searchButton = document.getElementById("search-button")
const searchInput = document.getElementById("search-input")
const resultsDiv = document.getElementById("results")

checkIfLoggedIn() 


searchButton.addEventListener("click", () => {

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

    searchResults.forEach(searchResult => {
        const item = document.createElement("div")
        item.classList.add("result-item")

        const link = document.createElement("a")
        link.href = searchResult.URL

        const titleSpan = document.createElement("span")
        titleSpan.classList.add("result-title")
        titleSpan.textContent = searchResult.title

        link.appendChild(titleSpan)
        item.appendChild(link)
        resultsDiv.appendChild(item)
    })

}





