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
    label.innerHTML = "Address"

    const line_break = document.createElement("br")

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
    row.appendChild(line_break)
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

