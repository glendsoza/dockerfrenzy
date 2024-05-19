import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import Snackbar from '@mui/material/Snackbar';
import MuiAlert from '@mui/material/Alert';
import useWebSocket from 'react-use-websocket';
import ReactJson from 'react-json-view';
import CircularProgress from '@mui/material/CircularProgress';
import {
  FileCopyOutlined as LogsIcon,
  StopOutlined as StopIcon,
  PlayArrowOutlined as StartIcon,
  RotateLeftOutlined as RestartIcon,
  LaunchOutlined as ExecIcon
} from '@mui/icons-material';
import NavBar from './NavBar'


function useQuery() {
  const { search } = useLocation();
  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function Container() {
  const [containerData, setContainerData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');
  const [snackbarSeverity, setSnackbarSeverity] = useState('success');
  const [buttonsDisabled, setButtonsDisabled] = useState(false); // State to control button disabled state

  const query = useQuery();

  useEffect(() => {
    setLoading(containerData.length === 0);
  }, [containerData]);

  const handleMessage = (event) => {
    setContainerData(JSON.parse(event.data)[0]);
  };

  const options = {
    onMessage: handleMessage
  };

  useWebSocket(`ws://${process.env.REACT_APP_API_SERVER_URL}/container/stream?ip=${query.get("ip")}&containerID=${query.get("containerID")}`, options);

  const performAction = (action) => {
    setButtonsDisabled(true); // Disable buttons when performing action
    fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/container/action?action=${action}&ip=${query.get("ip")}&containerID=${query.get("containerID")}`)
      .then(resp => resp.json())
      .then(data => {
        if (data.Error) {
          setSnackbarSeverity('error');
          setSnackbarMessage(data.Error);
        } else {
          setSnackbarSeverity('success');
          setSnackbarMessage(`Container ${action}ed successfully`);
        }
        setSnackbarOpen(true);
        setButtonsDisabled(false); // Enable buttons after action is completed
      });
  };

  const handleLogs = () => {
    window.open(`/container/log?ip=${query.get("ip")}&containerID=${query.get("containerID")}`, '_blank');
  };

  const execIntoContainer = () => {
    window.open(`/container/exec?ip=${query.get("ip")}&containerID=${query.get("containerID")}`, '_blank');
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbarOpen(false);
  };

  return (
    <div>
      <NavBar linkMap={[{ link: "/", name: "Machines" }, { link: `/machine?ip=${query.get("ip")}`, name: "Machine Info" }, { link: `http://${process.env.REACT_APP_EXEC_SERVER_URL}/config/show`, name: "Config" }]} />
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <div style={{
          backgroundColor: 'rgb(249, 245, 227)', // Changed background color to lighter
          border: '1px solid #ddd',
          borderRadius: '8px',
          padding: '20px',
          marginBottom: '20px',
          width: '300px',
          boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)',
          minHeight: '200px',
        }}>
          <h2 style={{ textAlign: 'center', marginBottom: '20px', fontFamily: 'Arial', color: '#333' }}>Container Details</h2>
          {!loading && containerData && (
            <div style={{ textAlign: 'left', fontFamily: 'Arial', color: '#333' }}>
              <p><b>Status:</b> {containerData.State.Status}</p>
              <p><b>Name:</b> {containerData.Name.replace(/^\/+/, '')}</p>
              <p><b>Restart Count:</b> {containerData.RestartCount}</p>
              <p><b>Image:</b> {containerData.Config.Image}</p>
              {/* Add more basic container info here */}
            </div>
          )}
          {loading && (
            <div style={{ textAlign: 'center' }}>
              <CircularProgress />
            </div>
          )}

          <div style={{ display: 'flex', justifyContent: 'space-around', marginTop: '20px' }}>
            <Tooltip title="View Logs" arrow>
              <IconButton onClick={handleLogs} disabled={loading || buttonsDisabled} style={{ backgroundColor: '#fff', boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)', borderRadius: '4px' }}>
                <LogsIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Stop Container" arrow>
              <IconButton onClick={() => performAction('stop')} disabled={loading || buttonsDisabled} style={{ backgroundColor: '#fff', boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)', borderRadius: '4px' }}>
                <StopIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Start Container" arrow>
              <IconButton onClick={() => performAction('start')} disabled={loading || buttonsDisabled} style={{ backgroundColor: '#fff', boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)', borderRadius: '4px' }}>
                <StartIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Restart Container" arrow>
              <IconButton onClick={() => performAction('restart')} disabled={loading || buttonsDisabled} style={{ backgroundColor: '#fff', boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)', borderRadius: '4px' }}>
                <RestartIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Exec Into Container" arrow>
              <IconButton onClick={execIntoContainer} disabled={loading || buttonsDisabled} style={{ backgroundColor: '#fff', boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)', borderRadius: '4px' }}>
                <ExecIcon />
              </IconButton>
            </Tooltip>
          </div>
        </div>
      </div>
      <script>

      </script>

      {!loading && containerData && (
        <ReactJson src={containerData} theme="monokai" />
      )}

      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <MuiAlert elevation={6} variant="filled" onClose={handleCloseSnackbar} severity={snackbarSeverity}>
          {snackbarMessage}
        </MuiAlert>
      </Snackbar>
    </div>
  );
}
