import './App.css';
import React from 'react';
import ReactDOM from 'react-dom';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid'
import Box from '@material-ui/core/Box'
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core'
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import lvicon from './lvicon.png';
import LockIcon from '@material-ui/icons/Lock';
import LockOpenIcon from '@material-ui/icons/LockOpen';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(25, 2),
    height: 400,
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center"
  },
  logo: {
    maxWidth: "2%",
    marginRight: theme.spacing(2),
  },
}));

function NavBar() {
  const classes = useStyles();
  return (
    <AppBar position="fixed" color="dark">
      <Toolbar>
        <img src={lvicon} alt="logo" className={classes.logo} />
        <Typography variant="h8">
          Stock Tracker
        </Typography>
      </Toolbar>
    </AppBar>
  );
}

class LoginForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = { isLoggedIn: true };

    this.handleClick = this.handleClick.bind(this);
  }

  handleClick() {
    this.setState(state => ({
      isLoggedIn: !state.isLoggedIn
    }));
  }

  render() {
    return (
      <div>
      <Grid
        container
        direction="row"
        justify="center"
        alignItems="center"
        spacing={2}
      >
        <Grid item xs={12}>
          {this.state.isLoggedIn? <LockIcon></LockIcon> : <LockOpenIcon></LockOpenIcon>}
        </Grid>
        <Grid item xs={2}></Grid>
        <Grid item xs={8}>
          <TextField id="standard-basic" label="Email Address" required fullWidth></TextField>
        </Grid>
        <Grid item xs={2}></Grid>
        <Grid item xs={2}></Grid>
        <Grid item xs={8}>
          <TextField id="standard-password-input" label="Password" type="password" required fullWidth></TextField>
        </Grid>
        <Grid item xs={2}></Grid>
        <Grid item xs={12}>
          <Button variant="contained" color="primary" onClick={this.handleClick}>Login</Button>
        </Grid>
      </Grid>
    </div>
    );
  }
}

function App() {
  const classes = useStyles();
  return (
    <div className="App">
      <NavBar></NavBar>
      <Box className={classes.root}>
        <Card>
          <CardContent>
            <LoginForm/>
          </CardContent>
          <CardActions>
          </CardActions>
        </Card>
      </Box>
    </div>
  );
}

export default App;
