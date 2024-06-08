import { Fragment } from "react";
import { Routes, Route } from "react-router-dom";
import { SnackbarProvider, closeSnackbar } from 'notistack'
import Machines from './components/Machines'
import MachineInfo from './components/MachineInfo'
import Container from './components/Container'
import Image from './components/Image'
import MachineExec from "./components/MachineExec";
import ContainerExec from "./components/ContainerExec";
import ContainerLog from "./components/ContainerLog";
import Config from "./components/Config";
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';

const action = snackbarId => (
  <Fragment>
    <IconButton
      size="small"
      aria-label="close"
      color="inherit"
      onClick={() => { closeSnackbar(snackbarId) }}
    >
      <CloseIcon fontSize="small" />
    </IconButton>
  </Fragment>
);

export default function App() {
  return (
    <div>
      <SnackbarProvider autoHideDuration={null} action={action} anchorOrigin={{ vertical: "top", horizontal: "right", }} preventDuplicate={true} style={{ "marginTop": "20px" }}>
        <Routes>
          <Route path="/" element={<Machines />} />
          <Route path="/machine" element={<MachineInfo />} />
          <Route path="/container" element={<Container />} />
          <Route path="/image" element={<Image />} />
          <Route path="/machine/exec" element={<MachineExec />} />
          <Route path="/container/exec" element={<ContainerExec />} />
          <Route path="/container/log" element={<ContainerLog />} />
          <Route path="/config" element={<Config />} />
        </Routes>
      </SnackbarProvider >
    </div>
  );
}

