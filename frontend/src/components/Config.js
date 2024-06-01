import React, { useEffect, useState } from "react";
import { linter, lintGutter } from "@codemirror/lint";
import CodeMirror from "@uiw/react-codemirror";
import * as yamlMode from "@codemirror/legacy-modes/mode/yaml";
import { StreamLanguage } from "@codemirror/language";
import { enqueueSnackbar } from "notistack";
import parser from "js-yaml";
import { githubLight } from "@uiw/codemirror-theme-github";
import NavBar from "./NavBar";
import IconButton from '@mui/material/IconButton';
import PublishIcon from '@mui/icons-material/Publish';
import Tooltip from '@mui/material/Tooltip';
import CachedIcon from '@mui/icons-material/Cached';

const yaml = StreamLanguage.define(yamlMode.yaml);

const yamlLinter = linter((view) => {
    const diagnostics = [];

    try {
        parser.load(view.state.doc);
    } catch (e) {
        const loc = e.mark;
        const from = loc ? loc.position : 0;
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

export default function Config(props) {
    const [config, setConfig] = useState("")
    const [reloadButtonDisabled, setReloadButtonDisabled] = useState(false)
    useEffect(() => {
        setReloadButtonDisabled(true)
        fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/config`).
            then(resp => resp.json()).
            then(data => {
                setConfig(data.Data)
                setReloadButtonDisabled(false)
            }).
            catch(error => {
                console.log(error)
                enqueueSnackbar("something went wrong while connecting to api server", { variant: "error" })
                setReloadButtonDisabled(false)
            })
    }, [])
    const onChange = React.useCallback((val, viewUpdate) => {
        setConfig(val);
    }, []);

    const handleReloadConfig = (event) => {
        setReloadButtonDisabled(true)
        fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/config/reload`,
            {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
            }).
            then(resp => resp.json()).
            then(data => {
                if (data.Error !== "") {
                    enqueueSnackbar(data.Error, { variant: "error" })
                } else {
                    enqueueSnackbar("config reloaded", { variant: "success" })
                }
                setReloadButtonDisabled(false)
            }).
            catch(error => {
                console.log(error)
                enqueueSnackbar("something went wrong while connecting to api server", { variant: "error" })
                setReloadButtonDisabled(false)
            })
    }
    return (
        <div>
            <NavBar linkMap={[{ link: "/", name: "Machines" }]} />
            <Tooltip title="Update Config" arrow>
                <IconButton onClick={handleReloadConfig} disabled={reloadButtonDisabled} style={{ backgroundColor: '#fff', borderRadius: '4px' }}>
                    <PublishIcon />
                </IconButton>
            </Tooltip>
            <Tooltip title="Reload Config" arrow>
                <IconButton onClick={handleReloadConfig} disabled={reloadButtonDisabled} style={{ backgroundColor: '#fff', borderRadius: '4px' }}>
                    <CachedIcon />
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
    )

}


