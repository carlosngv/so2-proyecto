// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {runtime} from '../models';

export function BlockAllDevices():Promise<void>;

export function BlockUSBPorts():Promise<void>;

export function DiskUsage():Promise<main.DiskStatus>;

export function GetCPUPercentage():Promise<main.CPUData>;

export function Greet(arg1:string):Promise<string>;

export function LogInfo(arg1:string):Promise<string>;

export function ReadMemoryStats():Promise<main.Memory>;

export function UnblockAllDevices():Promise<void>;

export function UnblockUSBPorts():Promise<void>;

export function WailsInit(arg1:runtime.Runtime):Promise<void>;
