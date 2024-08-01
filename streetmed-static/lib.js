
function submit(event) {
    const formdata = { q: event.target.value };
    const queryString = (new URLSearchParams(formdata)).toString();
    const response = document.getElementById("response")
    
    fetch(`/streetmed/api/protocol?${queryString}`)
        .then(resp => Promise.all([Promise.resolve(resp.ok), resp.text()]))
        .then(arr => {
            const ok = arr[0];
            const body = arr[1];
            response.innerHTML = body;

        })
        .catch(err => console.error(err));;

    return false;
}

