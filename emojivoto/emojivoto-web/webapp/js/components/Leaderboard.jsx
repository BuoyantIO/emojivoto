import React from 'react';
import _ from 'lodash';
import 'whatwg-fetch';
import { Link } from 'react-router-dom';


// TODO: Remove? (not currently)
// const randomLeaderboard = () => _.map(emo => {
//   return {
//     unicode: emo.unicode,
//     votes: Math.round(Math.random() * 1000)
//   }
// });
export default class Leaderboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      leaderboard: []
    }
  }

  componentDidMount() {
    this.loadFromServer();
  }

  loadFromServer(emoji) {
    fetch('/api/leaderboard').then(r => {
      r.json().then(emojis => {
        this.setState({
          leaderboard: emojis
        });
      }).catch(e => this.setState({ error: e }));
    }).catch(e => this.setState({ error: e }));

    // TODO: Remove? (not currently)
    // let leaderboard = randomLeaderboard();
    // this.setState({ leaderboard: _.orderBy(leaderboard, 'votes', 'desc') });
  }

  renderLeaderboard() {
    return _.map(this.state.leaderboard, (emoji, i) => {
      return (
        <div className="emoji" key={`emoji-${i}`} title={`${emoji.votes} votes`}>
          <div>{emoji.unicode}</div>
          { emoji.votes > 0 ? <div className="counter">{emoji.votes}</div> : null}
        </div>
      );
    });
  }

  render() {
    return (
      <div className="background">
        <div className="page-content">
            {!this.state.error ? null : <div className="error">Error loading leaderboard.</div>}
            <h1>EMOJI VOTE LEADERBOARD </h1>
            <div className="btn btn-blue"><Link to="/">Vote on your favorite</Link></div>
            <div className="emoji-list">{this.renderLeaderboard()}</div>
            <div className="conduit-footer">
              <p className="footer-cta">Tap here to learn more about Conduit</p>
              <p className="footer-cta-web">Click here to learn more about Conduit</p>
            </div>
        </div>
      </div>
    );
  }
}