//shortcut for macos
window.addEventListener("keypress", (event) => {  
    if (event.metaKey && event.key === 'c') {
        document.execCommand("copy")
        event.preventDefault();
    }
    if (event.metaKey && event.key === 'v') {
        document.execCommand("paste")
        event.preventDefault();
    }        
    if (event.metaKey && event.key === 'x') {
        document.execCommand("cut")
        event.preventDefault();
    }
    if (event.metaKey && event.key === 'z') {
        document.execCommand("undo")
        event.preventDefault();
    }
    if (event.metaKey && event.key === 'a') {
        document.execCommand("selectAll")
        event.preventDefault();
    }
})