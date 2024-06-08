import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { AgGridReact } from 'ag-grid-react';
import { enqueueSnackbar } from 'notistack'
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import { LaunchOutlined as ExecIcon } from '@mui/icons-material';
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-quartz.css';
import { FaCheck, FaTimes, FaUndo } from 'react-icons/fa';
import NavBar from './NavBar'
import './styles/GridTheme.css';
import './styles/ErrorModal.css'

const sshIntoMachine = () => {
  return (<Tooltip title="Exec Into Machine" arrow>
    <IconButton onClick={(event) => {
      event.preventDefault()
    }} style={{ borderRadius: '4px' }}>
      <ExecIcon />
    </IconButton>
  </Tooltip>
  )
}

const LoadingIndicator = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
    <div className="spinner-border" role="status">
      <span className="sr-only">Loading...</span>
    </div>
  </div>
);

const Modal = ({ onClose }) => (
  <div className="modal">
    <div className="modal-content">
      <span className="close" onClick={onClose}>&times;</span>
      <p>Sorry, this machine is currently offline and cannot be accessed.</p>
    </div>
  </div>
);


export default function Machines() {
  const navigate = useNavigate();
  const [rows, setRows] = useState([]);
  const [gridApi, setGridApi] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/machines`)
      .then(response => response.json())
      .then(data => {
        if (data.error !== "") {
          setRows(data.Machines);
          setIsLoading(false);
        } else {
          setRows([])
          enqueueSnackbar(data.Error, { variant: "error" })
          setIsLoading(false);
        }
      })
      .catch(error => {
        console.error(error)
        setRows([])
        enqueueSnackbar("something went wrong while connecting to api server", { variant: "error" })
        setIsLoading(false);
      });
  }, []);

  useEffect(() => {
    if (gridApi) {
      gridApi.setQuickFilter('');
    }
  }, [gridApi]);

  const handleRowSelection = (row) => {
    if (row.data.Status !== 1) {
      setShowModal(true); // Show modal if machine is offline
    } else if (row.event.defaultPrevented) {
      // then we navigate to machine exec
      window.open(`/machine/exec?ip=${row.data.Ip}`, '_blank');
    } else {
      navigate(`/machine?ip=${row.data.Ip}`);
    }
  };

  function clearFilters() {
    gridApi.setFilterModel(null);
  }

  const getStatusIcon = (params) => {
    if (params.value === 0) {
      return <FaTimes style={{ color: 'red', fontSize: '20px' }} />;
    } else if (params.value === 1) {
      return <FaCheck style={{ color: 'green', fontSize: '20px' }} />;
    }
    return '';
  };

  const statusFilterParams = {
    buttons: ['apply', 'reset'],
    closeOnApply: true,
    suppressAndOrCondition: true,
    textFormatter: (value) => value === 0 ? '✘' : '✔️',
    debounceMs: 200
  };

  const gridOptions = {
    autoSizeStrategy: { type: 'fitGridWidth' },
    resizable: true,
    sortable: true,
    wrapText: true,
    autoHeight: true,
    pagination: true,
    onRowClicked: handleRowSelection,
    // Add getRowClass to apply selectable cursor to entire row
    getRowClass: () => 'selectable-row',
    onGridReady: (params) => setGridApi(params.api),
    suppressRowClickSelection: true,
  };

  if (isLoading) {
    return <LoadingIndicator />;
  }

  return (
    <div>
      {showModal && <Modal onClose={() => setShowModal(false)} />}
      <NavBar linkMap={[{ link: `/config`, name: "Config" }]} />
      <h2 style={{ textAlign: 'center', margin: '20px 0' }}>Machines</h2>
      <div style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: '10px' }}>
        <button onClick={clearFilters} style={{ display: 'flex', alignItems: 'center', backgroundColor: '#007bff', color: '#fff', border: 'none', borderRadius: '8px', padding: '10px 20px', cursor: 'pointer', boxShadow: '0px 4px 6px rgba(0, 0, 0, 0.1)', transition: 'background-color 0.3s', }}>
          <FaUndo style={{ marginRight: '5px' }} />
          Reset Filters
        </button>
      </div>
      <div className="ag-theme-quartz" style={{ height: '80vh' }}>
        <AgGridReact
          gridOptions={gridOptions}
          columnDefs={[
            { headerName: 'HostName', field: 'HostName', filter: 'agTextColumnFilter', sortable: true },
            { headerName: 'Ip', tooltipField: 'Ip', field: 'Ip', filter: 'agTextColumnFilter', sortable: true },
            { headerName: 'Os', field: 'Os', filter: 'agTextColumnFilter', sortable: true },
            { headerName: 'Status', field: 'Status', cellRenderer: getStatusIcon, filter: 'agSetColumnFilter', filterParams: statusFilterParams, sortable: true },
            { headerName: 'Error', field: 'Error', tooltipField: 'Error', filter: 'agTextColumnFilter', sortable: true },
            { headerName: "ssh", field: "SSH into machine", cellRenderer: sshIntoMachine },
          ]}
          rowData={rows}
        ></AgGridReact>
      </div>
    </div>
  );
}
