let address_counter = 1;

function addAddressInput(event) {
    event.preventDefault();

    let new_id = address_counter + 1

    const submit_btn = document.getElementById("submit-btn");

    const form = document.getElementById("address-input-form");

    // Get row element with incrementing counter
    const row = document.createElement("div");
    row.setAttribute("id", "address-line-" + new_id);

    const label = document.createElement("label")
    label.setAttribute("for", "address_input_"+ new_id)

    const h3 = document.createElement("h3") 
    h3.innerHTML = "Address"

    const input = document.createElement("input")
    input.setAttribute("id", "address_input_" + new_id)
    input.setAttribute("type", "text")
    input.setAttribute("name", "address_" + new_id)
    input.className = "address-input"

    const button = document.createElement("button")
    button.setAttribute("onClick", "addAddressInput(event)")
    button.className = "new-address-btn"
    button.innerHTML = "+"

    // Insert row at end before submit button
    form.insertBefore(row, submit_btn)
    // Append everything to the row
    row.appendChild(label)
    label.appendChild(h3)
    row.appendChild(input)
    row.appendChild(button)
    row.appendChild(document.createElement("br"))

    // Change previous button to be a "-" instead of a "+" to allow users to delete rows
    let previous_row = document.getElementById("address-line-" + address_counter)
    let previous_button = previous_row.getElementsByClassName("new-address-btn")[0]

    previous_button.innerText = "-"
    previous_button.setAttribute("onClick", "removeAddressInput(event)")

    // Increment counter
    address_counter += 1

}

function removeAddressInput(event) {
    let button = (event.target) ? event.target: event.srcElement;
    let addressInput = button.closest("div")
    
    button.remove()
    addressInput.remove()
}



let address_layer;
let qp_layer;


// Function to give the submit button a waiting animation
function startWaitingAnimation() {
    let submit_btn = document.getElementById("submit-btn")
    submit_btn.setAttribute("aria-busy", "true")
    submit_btn.innerHTML = "Please wait..."
}

function endWaitingAnimation() {
    let submit_btn = document.getElementById("submit-btn")
    submit_btn.setAttribute("aria-busy", "false")
    submit_btn.innerHTML = "Submit"
}

let search_results;

function submitAddresses(e) {
    e.preventDefault();
    try {
        emptyTable();
    } catch(err) {
        console.log("No table exists")
    }

    // Start waiting animation
    startWaitingAnimation()

    // let content = document.getElementById("text-content")
    var form = document.getElementById("address-input-form");
    var formData = new FormData(form);

    let object = {}
    formData.forEach(function(value, key) {
        object[key] = value
    })

    let jsonFormData = JSON.stringify(object)

    fetch("/submit", {
        headers: {"Content-Type": "application/json"},
        method: "POST",
        body: jsonFormData,
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('network returns error');
        }
    return response.json();
    })
    .then((resp) => {
        // Remove all layers except the base
        try {
            map.eachLayer(function (layer) {
                if (layer != titleLayer) {
                    map.removeLayer(layer);
                }
            });
        }
        catch(err) {
            console.log("No layers exist")
        }
        // Map addresses
        let addressArray = resp["addresses"] 
        let addressDict = {} 

        for (let i = 0; i < addressArray.length; i++) {
            addressDict["address_" + i] = L.marker([addressArray[i][0], addressArray[i][1]])
        }
        let addressfeatureGroup = Object.values(addressDict)
        let address_layer = new L.featureGroup(addressfeatureGroup).addTo(map)

        // Map query points
        let qpArray = resp["query_points"] 
        let qpDict = {} 

        for (let i = 0; i < qpArray.length; i++) {
            qpDict["address_" + i] = L.circleMarker([qpArray[i][0], qpArray[i][1]])
        }
        let qpfeatureGroup = Object.values(qpDict)
        let qp_layer = new L.featureGroup(qpfeatureGroup).addTo(map)

        // Process search results and put into table
        search_results = resp["results"]
        // Only if yelp results were requested
        let map_container = document.getElementById("map-container")
        if ("results" in resp) {
            constructTable(search_results)
            map_container.className = "map-container-split"
            setTimeout(function(){ map.invalidateSize()}, 400);
        } else {
            map_container.className = "map-container-full"
            setTimeout(function(){ map.invalidateSize()}, 400);
        }
        // End waiting animation
        endWaitingAnimation()
        scrollTo("map_and results")
        map.fitBounds(address_layer.getBounds().pad(0.5))

    })
    .catch((error) => {
        // Handle error
        console.log("error ", error);
    });
}

var addressForm = document.getElementById("address-input-form");
addressForm.addEventListener("submit", submitAddresses);

function emptyTable() {
    let table = document.getElementById("results")
    let header = document.getElementById("table_header")
    let table_body = document.getElementById("table_body")

    table.removeChild(header)
    table.removeChild(table_body)
}

function scrollTo(hash) {
    location.hash = "#" + hash;
}

// Construct results table
function constructTable(json_data) {
    let table = document.getElementById("results")
    // Getting the all column names
    let col_names = ["Name", "Address", "Price", "URL"];

    // Create table header
    let thead = document.createElement("thead")
    thead.setAttribute("style", "position: sticky;") 
    let thead_tr = document.createElement("tr") 
    thead.setAttribute("id", "table_header")
    thead.appendChild(thead_tr)

    col_names.forEach(function(col) {
        let th = document.createElement("th")
        th.setAttribute("scope", "col")
        th.innerHTML = col 

        thead_tr.appendChild(th)
    })

    let tbody = document.createElement("tbody")
    // Traversing the JSON data
    for (let i = 0; i < json_data.length; i++) {
        let row = document.createElement("tr");

        // On hover of a row, map relevant address and name
        // Adding event listeners to rows
        row.id = "result_row_" + i
        row.addEventListener('mouseover', (obj) => {
            try {
                map.removeLayer(selected_marker)
            } catch {
                console.log("popup doesn't exist")
            }

            let relevant_name = search_results[i]["name"]
            let relevant_coords = search_results[i]["coordinates"]

            selected_marker = L.popup(closeButton=false).setLatLng(relevant_coords)
                .setContent(relevant_name).openOn(map)

        })
        row.addEventListener('click', (obj) => {
            try {
                map.removeLayer(selected_marker)
            } catch {
                console.log("popup doesn't exist")
            }

            let relevant_name = search_results[i]["name"]
            let relevant_coords = search_results[i]["coordinates"]

            selected_marker = L.popup(closeButton=false).setLatLng(relevant_coords)
                .setContent(relevant_name).openOn(map)

        })
        // row.addEventListener('mouseout', (obj) => {
        //     map.removeLayer(selected_marker)
        // })

        // Build table body
        for (let colIndex = 0; colIndex < col_names.length; colIndex++) {
            let val = json_data[i][col_names[colIndex].toLowerCase()];

            // If there is no key matching the column name
            if (val == null) val = "";

            let td = document.createElement("td")

            if (col_names[colIndex] === "URL") {
                let url_link = document.createElement("a")
                url_link.setAttribute("target", "_blank")
                url_link.href = val
                url_link.innerHTML = "Yelp"

                td.append(url_link)
            } else {
                td.innerHTML = val
            }
            
            row.append(td);
        }

        // Adding each row to the table
        tbody.append(row);
    }
    tbody.setAttribute("id", "table_body")

    table.appendChild(thead)
    table.appendChild(tbody)
}
