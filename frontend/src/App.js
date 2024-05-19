import { Routes, Route } from "react-router-dom";
import Machines from './components/Machines'
import MachineInfo from './components/MachineInfo'
import Container from './components/Container'
import Image from './components/Image'
import MachineExec from "./components/MachineExec";
import ContainerExec from "./components/ContainerExec";
import ContainerLog from "./components/ContainerLog";

export default function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Machines />} />
        <Route path="/machine" element={<MachineInfo />} />
        <Route path="/container" element={<Container />} />
        <Route path="/image" element={<Image />} />
        <Route path="/machine/exec" element={<MachineExec />} />
        <Route path="/container/exec" element={<ContainerExec />} />
        <Route path="/container/log" element={<ContainerLog />} />
      </Routes>
    </div>
  );
}

