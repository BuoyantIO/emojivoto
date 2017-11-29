import React from 'react';
import _ from 'lodash';
import { Link } from 'react-router-dom';
import 'whatwg-fetch';
export default class Vote extends React.Component {
  constructor(props) {
    super(props);
    this.resetState = this.resetState.bind(this);
    this.state = {
      emojiList: [],
      selectedEmoji: null,
      error: null
    }
  }

  loadFromServer() {
    fetch('/api/list').then(rsp => {
      rsp.json().then(r => {
        this.setState({ emojiList: r })
      }).catch(e => {
        this.setState( {error: e.statusText });
      }).catch(e => {
        this.setState( {error: e.statusText });
      });
    });
  }

  componentDidMount() {
    this.loadFromServer();
  }

  vote(emoji) {
    fetch(`/api/vote?choice=${emoji.shortcode}`)
      .then(rsp => {
        if (rsp.ok) {
          this.setState({ selectedEmoji: emoji, error: null });
        } else {
          this.setState({ error: rsp.statusText });
        }
      })
      .catch(e => {
        this.setState({ error: e.statusText })
      });
    this.setState({ selectedEmoji: emoji }); // TODO: remove
  }

  resetState() {
    this.setState({ selectedEmoji: null, error: null });
  }

  renderEmojiList(emojis) {
    return _.map(emojis, (emoji, i) => {
      return (
        <div
          className="emoji"
          key={`emoji-${i}`}
          onClick={e => this.vote(emoji)}
        >
          {emoji.unicode}
        </div>
      );
    });
  }

  render() {
    if (!this.state.selectedEmoji) {
      let emojiList = this.state.emojiList;
      return (
        <div className="background">
          <div className="page-content" data-aos="fade-left">
            {!this.state.error ? null :
              <div className="error">Could not vote. Error: {this.state.error}</div>}
            <h1 className="headline">🗳</h1>
            <h1>VOTEMOJI</h1>
            <p>Tap to vote for your favorite emoji below</p>
            <div className="btn btn-blue"><Link to="/leaderboard">View the leaderboard</Link></div>
            {!_.isEmpty(this.state.emojiList) ? null : <div>Loading emoji...</div>}
            <div className="emoji-list">{this.renderEmojiList(emojiList)}</div>
          </div>
        </div>
      );
    } else {
      return (
        <div className="background">
          <div className="page-content">
            <h1>You picked:</h1>
            <h1 className="headline">{this.state.selectedEmoji.unicode}</h1>
            <p>See how you stack up against others</p>
            <div className="btn btn-blue"><Link to="/leaderboard">View the leaderboard</Link></div>
            <div className="btn btn-white"><Link to="/" onClick={this.resetState}>Pick another one</Link></div>
          </div>
        </div>
      );
    }
  }
}