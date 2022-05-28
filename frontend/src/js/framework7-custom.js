
import Framework7, { request, utils, getDevice, createStore } from 'framework7';
import Calendar from 'framework7/components/calendar';
import Card from 'framework7/components/card';
import Checkbox from 'framework7/components/checkbox';
import ContactsList from 'framework7/components/contacts-list';
import Chip from 'framework7/components/chip';
import DataTable from 'framework7/components/data-table';
import Dialog from 'framework7/components/dialog';
import Form from 'framework7/components/form';
import Grid from 'framework7/components/grid';
import InfiniteScroll from 'framework7/components/infinite-scroll';
import Input from 'framework7/components/input';
import LoginScreen from 'framework7/components/login-screen';
import Panel from 'framework7/components/panel';
import Popover from 'framework7/components/popover';
import Popup from 'framework7/components/popup';
import Preloader from 'framework7/components/preloader';
import Progressbar from 'framework7/components/progressbar';
import Radio from 'framework7/components/radio';
import Searchbar from 'framework7/components/searchbar';
import Skeleton from 'framework7/components/skeleton';
import Toast from 'framework7/components/toast';
import Toggle from 'framework7/components/toggle';
import Tooltip from 'framework7/components/tooltip';
import Typography from 'framework7/components/typography';
import VirtualList from 'framework7/components/virtual-list';
import Fab from 'framework7/components/fab';

Framework7.use([
  Calendar,
  Card,
  Checkbox,
  ContactsList,
  Chip,
  DataTable,
  Dialog,
  Form,
  Grid,
  InfiniteScroll,
  Input,
  LoginScreen,
  Panel,
  Popover,
  Popup,
  Preloader,
  Progressbar,
  Radio,
  Searchbar,
  Skeleton,
  Toast,
  Toggle,
  Tooltip,
  Typography,
  VirtualList,
  Fab,
]);

export default Framework7;
export { request, utils, getDevice, createStore };
