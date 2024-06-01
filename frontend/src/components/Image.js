import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import useWebSocket from 'react-use-websocket';
import JsonView from 'react18-json-view';
import 'react18-json-view/src/style.css'
import CircularProgress from '@mui/material/CircularProgress';
import {
  LaunchOutlined as CreateContainerIcon
} from '@mui/icons-material';
import NavBar from './NavBar';
import CreateContainer from './CreateContainer'

function useQuery() {
  const { search } = useLocation();
  return React.useMemo(() => new URLSearchParams(search), [search]);
}

function bytesToSize(bytes) {
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  if (bytes === 0) return '0 Byte';
  const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)), 10);
  return `${Math.round(bytes / Math.pow(1024, i), 2)} ${sizes[i]}`;
}

export default function Container() {
  const [imageData, setImageData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);

  const query = useQuery();

  useEffect(() => {
    setLoading(imageData.length === 0);
  }, [imageData]);

  const handleMessage = (event) => {
    setImageData(JSON.parse(event.data)[0]);
  };

  const options = {
    onMessage: handleMessage
  };

  useWebSocket(`ws://${process.env.REACT_APP_API_SERVER_URL}/image/stream?ip=${query.get("ip")}&imageID=${query.get("imageID")}`, options);

  const handleCreateContainer = () => {
    setModalOpen(true);
  };

  const handleCloseModal = () => {
    setModalOpen(false);
  };



  return (
    <div>
      <NavBar linkMap={[{ link: "/", name: "Machines" }, { link: `/machine?ip=${query.get("ip")}`, name: "Machine Info" }, { link: `/config`, name: "Config" }]} />
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <div style={{
          backgroundColor: 'rgb(249, 245, 227)',
          border: '1px solid #ddd',
          borderRadius: '8px',
          padding: '20px',
          marginBottom: '20px',
          width: '300px',
          boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)',
          minHeight: '200px',
          overflowX: 'auto', // Enable horizontal scrollbar
          WebkitOverflowScrolling: 'touch', // Mobile support
          // Customize scrollbar appearance
          scrollbarWidth: 'thin',
          scrollbarColor: '#ccc #f0f0f0',
        }}>
          <h2 style={{ textAlign: 'center', marginBottom: '20px', fontFamily: 'Arial', color: '#333' }}>Image Details</h2>
          {!loading && imageData && (
            <div style={{ textAlign: 'left', fontFamily: 'Arial', color: '#333' }}>
              <p><b>Architecture:</b> {imageData.Architecture}</p>
              <p><b>Os:</b> {imageData.Os}</p>
              <p><b>Size:</b> {bytesToSize(imageData.Size)}</p>
            </div>
          )}
          {loading && (
            <div style={{ textAlign: 'center' }}>
              <CircularProgress />
            </div>
          )}

          <div style={{ display: 'flex', justifyContent: 'space-around', marginTop: '20px' }}>
            <Tooltip title="Create Container" arrow>
              <IconButton onClick={handleCreateContainer} disabled={loading} style={{ backgroundColor: '#fff', boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)', borderRadius: '4px' }}>
                <CreateContainerIcon />
              </IconButton>
            </Tooltip>
          </div>
        </div>
      </div>

      {!loading && imageData && (
        <JsonView src={imageData} theme="vscode" />
      )}
      {/* Modal Form */}
      {modalOpen && (
        <div><CreateContainer handleCloseModal={handleCloseModal} /></div>
      )}
    </div>
  );
}
