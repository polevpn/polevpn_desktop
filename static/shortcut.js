//shortcut for macos
window.addEventListener("keypress", (event) => {  
    if (event.metaKey && event.key === 'c') {
        document.execCommand("copy",false)
        event.preventDefault();
    }
    if (event.metaKey && event.key === 'v') {
        document.execCommand("paste",false)
        event.preventDefault();
    }        
    if (event.metaKey && event.key === 'x') {
        document.execCommand("cut",false)
        event.preventDefault();
    }
    if (event.metaKey && event.key === 'z') {
        document.execCommand("undo",false)
        event.preventDefault();
    }
    if (event.metaKey && event.key === 'a') {
        document.execCommand("selectAll",false)
        event.preventDefault();
    }
})