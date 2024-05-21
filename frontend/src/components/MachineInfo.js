import React, { useState,useEffect,useCallback } from 'react';
import { useLocation } from 'react-router-dom';
import { useNavigate } from "react-router-dom";
import { AgGridReact } from 'ag-grid-react';
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-quartz.css';
import './styles/GridTheme.css';
import { FaUndo } from 'react-icons/fa';
import NavBar from './NavBar'

function useQuery() {
  const { search } = useLocation();
  return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function MachineInfo() {
  const navigate = useNavigate();
  let query = useQuery();
  const [containerGridApi, setContainerGridApi] = useState(null);
  const [imageGridApi, setImageGridApi] = useState(null);
  const [containerData, setContainerData] = useState(null);
  const [imageData, setImageData] = useState(null);

  useEffect(() => {
    if (containerGridApi) {
      containerGridApi.setQuickFilter('');
    }
  }, [containerGridApi]);

  useEffect(() => {
    if (imageGridApi) {
      imageGridApi.setQuickFilter('');
    }
  }, [imageGridApi]);

  function clearContainerFilters() {
    containerGridApi.setFilterModel(null);
  }

  function clearImageFilters() {
    imageGridApi.setFilterModel(null);
  }

  const handleContainerRowSelection = (row) => {
    navigate(`/container?ip=${query.get("ip")}&containerID=${row.data.ID}`)
  }
  const handleImageRowSelection = (row) => {
    navigate(`/image?ip=${query.get("ip")}&imageID=${row.data.ID}`)
  }

  const onContainerGridReady = useCallback((params) => {
    setContainerGridApi(params.api)
    fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/containers?ip=${query.get("ip")}`)
    .then(resp => resp.json())
    .then(data => {
      let containerData = data.Containers
      if (containerData === null) {
        containerData = []
      }
      setContainerData(containerData)
    })
  }
  )

  const onImageGridReady = useCallback((params) => {
    setImageGridApi(params.api)
    fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/images?ip=${query.get("ip")}`)
    .then(resp => resp.json())
    .then(data => {
      let imageData = data.Images
      if (imageData === null) {
        imageData = []
      }
      setImageData(imageData)
    })
  }
  )
  const containerGridOptions = {
    autoSizeStrategy: { type: 'fitGridWidth' },
    wrapText: true,
    autoHeight: true,
    pagination: true,
    onRowClicked: handleContainerRowSelection,
    getRowClass: () => 'selectable-row',
    onGridReady: onContainerGridReady, 
  };

  const imageGridOptions = {
    autoSizeStrategy: { type: 'fitGridWidth' },
    wrapText: true,
    autoHeight: true,
    pagination: true,
    onRowClicked: handleImageRowSelection,
    getRowClass: () => 'selectable-row',
    onGridReady: onImageGridReady,
  };
  return (
    <div>
      <NavBar linkMap={[{link: "/", name:"Machines"},{link:`/config`, name:"Config"}]}/>
    <div>
    <h2 style={{ textAlign: 'center'}}>Containers</h2>
    <div style={{ display: 'flex', justifyContent: 'flex-end'}}>
      <button onClick={clearContainerFilters} style={{ display: 'flex', alignItems: 'center', backgroundColor: '#007bff', color: '#fff', border: 'none', borderRadius: '8px', padding: '10px 20px', cursor: 'pointer', boxShadow: '0px 4px 6px rgba(0, 0, 0, 0.1)', transition: 'background-color 0.3s', }}>
        <FaUndo style={{ marginRight: '5px' }} />
        Reset Filters
      </button>
    </div>
    <div className="ag-theme-quartz" style={{ height: '50vh' }}>
      <AgGridReact
        gridOptions={containerGridOptions}
        columnDefs={[
          { headerName: 'Names', field: 'Names', filter: 'agTextColumnFilter', tooltipField: "Names"},
          { headerName: 'Image',field: 'Image', filter: 'agTextColumnFilter', sortable: true},
          { headerName: 'State', field: 'State', filter: 'agTextColumnFilter', sortable: true },
          { headerName: 'Ports', field: 'Ports',  filter: 'agTextColumnFilter'},
        ]}
        rowData={containerData}
      ></AgGridReact>
    </div>
  </div>

  <div>
    <h2 style={{ textAlign: 'center'}}>Images</h2>
    <div style={{ display: 'flex', justifyContent: 'flex-end'}}>
      <button onClick={clearImageFilters} style={{ display: 'flex', alignItems: 'center', backgroundColor: '#007bff', color: '#fff', border: 'none', borderRadius: '8px', padding: '10px 20px', cursor: 'pointer', boxShadow: '0px 4px 6px rgba(0, 0, 0, 0.1)', transition: 'background-color 0.3s', }}>
        <FaUndo style={{ marginRight: '5px' }} />
        Reset Filters
      </button>
    </div>
    <div className="ag-theme-quartz" style={{ height: '50vh' }}>
      <AgGridReact
        gridOptions={imageGridOptions}
        columnDefs={[
          { headerName: 'ID', field: 'ID', filter: 'agTextColumnFilter', tooltipField: "ID"},
          { headerName: 'Tag',field: 'Tag', filter: 'agTextColumnFilter', sortable: true},
          { headerName: 'Repository', field: 'Repository', filter: 'agTextColumnFilter', sortable: true,tooltipField: "Repository" },
          { headerName: 'Size', field: 'Size',  filter: 'agTextColumnFilter'},
          { headerName: 'Containers', field: 'Containers', filter: 'agTextColumnFilter', tooltipField: "Containers"},
        ]}
        rowData={imageData}
      ></AgGridReact>
    </div>
  </div>
  </div>
);
}