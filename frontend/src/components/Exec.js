import React, { useEffect } from 'react';
import "xterm/css/xterm.css"
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';


export default function Exec(props) {
    useEffect(() => {
        if (document.getElementsByClassName("terminal").length === 0) {
            let term = new Terminal()
            const fitAddon = new FitAddon();
            term.loadAddon(fitAddon);
            term.open(document.getElementById('terminal'));
            fitAddon.fit()
            var socket = new WebSocket(props.url)
            socket.addEventListener("message", async (event) => {
                if (event.data.constructor.name === 'Blob') {
                    term.write(await event.data.text())
                } else {
                    term.write(event.data)
                }
            })
            if (!props.readOnly) {
                term.onData((data) => socket.send(data))
            }
        }
    }, [])
    return (
        <div id="terminal" style={{ width: "100vw", height: "100vh" }}></div>
    )
}