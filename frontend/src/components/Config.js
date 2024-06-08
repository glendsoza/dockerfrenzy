import React, { useEffect, useState } from "react";
import { linter, lintGutter } from "@codemirror/lint";
import CodeMirror from "@uiw/react-codemirror";
import { StreamLanguage } from "@codemirror/language";
import { enqueueSnackbar } from "notistack";
import parser from "js-yaml";
import { githubLight } from "@uiw/codemirror-theme-github";
import NavBar from "./NavBar";
import IconButton from '@mui/material/IconButton';
import PublishIcon from '@mui/icons-material/Publish';
import Tooltip from '@mui/material/Tooltip';
import ComputerIcon from '@mui/icons-material/Computer';
import TextField from '@mui/material/TextField';

// Import YAML mode correctly
import { yaml as yamlMode } from "@codemirror/legacy-modes/mode/yaml";

// Define the YAML language
const yaml = StreamLanguage.define(yamlMode);

// Define the linter for YAML
const yamlLinter = linter((view) => {
    const diagnostics = [];
    const doc = view.state.doc.toString();

    try {
        parser.load(doc);
    } catch (e) {
        const loc = e.mark;
        const from = loc ? Math.min(loc.position, doc.length - 1) : 0;
        const to = from;
        const severity = "error";

        diagnostics.push({
            from,
            to,
            message: e.message,
            severity
        });
    }
    return diagnostics;
});

export default function Config() {
    const [config, setConfig] = useState("");
    const [buttonDisabled, setButtonDisabled] = useState(false);

    // Fetch the config data on component mount
    useEffect(() => {
        setButtonDisabled(true);
        fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/config`)
            .then(resp => resp.json())
            .then(data => {
                setConfig(data.Data);
                setButtonDisabled(false);
            })
            .catch(error => {
                console.error(error);
                enqueueSnackbar("Something went wrong while connecting to API server", { variant: "error" });
                setButtonDisabled(false);
            });
    }, []);

    // Handle change in CodeMirror editor
    const onChange = React.useCallback((val, viewUpdate) => {
        setConfig(val);
    }, []);

    // Handle button press actions
    const handleButtonPress = (event, url, action, payload) => {
        setButtonDisabled(true);
        fetch(url, {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ config: payload })
        })
            .then(resp => resp.json())
            .then(data => {
                if (data.Error !== "") {
                    enqueueSnackbar(data.Error, { variant: "error" });
                } else {
                    enqueueSnackbar(`Config ${action}`, { variant: "success" });
                }
                setButtonDisabled(false);
            })
            .catch(error => {
                console.error(error);
                enqueueSnackbar("Something went wrong while connecting to API server", { variant: "error" });
                setButtonDisabled(false);
            });
    };

    // Handle config update
    const handleUpdateConfig = (event) => {
        handleButtonPress(event, `http://${process.env.REACT_APP_API_SERVER_URL}/config/update`, "updated", config);
    };

    // Handle config reload
    const handleReloadConfig = (event) => {
        handleButtonPress(event, `http://${process.env.REACT_APP_API_SERVER_URL}/config/reload`, "reloaded", "");
    };

    return (
        <div>
            <NavBar linkMap={[{ link: "/", name: "Machines" }]} />
            <TextField
                disabled
                fullWidth={true}
                id="outlined-disabled"
                label="Example"
                multiline
                maxRows={Infinity}
                defaultValue={`# Private key file is always relative to the config folder env variable
# No comments allowed in the actual config file          
# configList:
#   - sshConfig:
#       passwordAuth:
#         username: glen
#         password: password
#       ips:
#         - 
#   - sshConfig:
#       sshAuth:
#         username: glen
#         PrivateKeyFile: password
#       ips:
#         - 
                `}
            />
            <Tooltip title="Update Config" arrow>
                <IconButton onClick={handleUpdateConfig} disabled={buttonDisabled} style={{ backgroundColor: '#fff', borderRadius: '4px' }}>
                    <PublishIcon />
                </IconButton>
            </Tooltip>
            <Tooltip title="Reload Machines" arrow>
                <IconButton onClick={handleReloadConfig} disabled={buttonDisabled} style={{ backgroundColor: '#fff', borderRadius: '4px' }}>
                    <ComputerIcon />
                </IconButton>
            </Tooltip>
            <CodeMirror
                height="100vh"
                value={config}
                theme={githubLight}
                extensions={[yaml, lintGutter(), yamlLinter]}
                onChange={onChange}
            />
        </div>
    );
}
