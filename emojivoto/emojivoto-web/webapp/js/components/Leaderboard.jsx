import React from 'react';
import _ from 'lodash';
import 'whatwg-fetch';
import { Link } from 'react-router-dom';

const randomLeaderboard = () => _.map(Emoji, emo => {
  return {
    unicode: emo.unicode,
    votes: Math.round(Math.random() * 1000)
  }
});
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

    let leaderboard = randomLeaderboard();
    this.setState({ leaderboard: _.orderBy(leaderboard, 'votes', 'desc') });
  }

  renderLeaderboard() {
    return _.map(this.state.leaderboard, (emoji, i) => {
      return (
        <div className="emoji" key={`emoji-${i}`} title={`${emoji.votes} votes`}>
          <div>{emoji.unicode}</div>
          { emoji.votes > 0 ? <div>{emoji.votes}</div> : null}
        </div>
      );
    });
  }

  render() {
    return (
      <div className="background">
        <div className="page-content">
            {!this.state.error ? null : <div className="error">Error loading leaderboard.</div>}
            <h1>VOTEMOJI LEADERBOARD </h1>
            <div className="btn btn-blue"><Link to="/">Vote on your favorite</Link></div>
            <div className="emoji-list">{this.renderLeaderboard()}</div>
        </div>
      </div>
    );
  }
}