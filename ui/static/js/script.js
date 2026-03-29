let IntervalID = setInterval(() => {
    fetch(window.location.href).then(response => response.json()).then(data => {
    if (data["WindowChanged"]){
        window.location.reload();
        }
    })
}, 50 * 1000)