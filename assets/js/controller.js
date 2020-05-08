const fill = 0xaaaaaa;
const line = 0x888888;
const active = 0.5;
const inactive = 0.3;

class ABXY extends PIXI.Graphics {
  constructor(app, x, y) {
    const abxy = super();

    const r = 30;
    this.input = { a: 0, b: 0, x: 0, y: 0 };
    this.config = [
      { x: (r * 5) / 2, y: r, label: "x" },
      { x: r, y: (r * 5) / 2, label: "y" },
      { x: r * 4, y: (r * 5) / 2, label: "a" },
      { x: (r * 5) / 2, y: r * 4, label: "b" }
    ];

    for (const cfg of this.config) {
      const button = new PIXI.Graphics();

      button.interactive = true;
      button.buttonMode = true;
      button.alpha = inactive;

      button
        .lineStyle(2, line)
        .beginFill(fill)
        .drawCircle(cfg.x, cfg.y, r)
        .endFill()
        .on("pointerdown", this.pointerDown.bind(this, cfg.label, button))
        .on("pointerup", this.pointerUp.bind(this, cfg.label, button))
        .on("pointerupoutside", this.pointerUp.bind(this, cfg.label, button));
      abxy.addChild(button);
    }

    app.stage.addChild(abxy);
    this.transform.position.x = x;
    this.transform.position.y = y;
  }

  pointerDown(label, button) {
    this.input[label] = 1;
    button.alpha = active;
  }

  pointerUp(label, button) {
    this.input[label] = 0;
    button.alpha = inactive;
  }
}

class Button extends PIXI.Graphics {
  constructor(app, label, x, y) {
    const button = super();

    this.input = { [label]: 0 };
    const r = 20;

    button.interactive = true;
    button.buttonMode = true;
    button.alpha = inactive;

    button
      .lineStyle(2, line)
      .beginFill(fill)
      .drawCircle(0, 0, r)
      .endFill()
      .on("pointerdown", this.pointerDown.bind(this, label))
      .on("pointerup", this.pointerUp.bind(this, label))
      .on("pointerupoutside", this.pointerUp.bind(this, label));

    app.stage.addChild(button);
    this.transform.position.x = x;
    this.transform.position.y = y;
  }

  pointerDown(label) {
    this.input[label] = 1;
    this.alpha = active;
  }

  pointerUp(label) {
    this.input[label] = 0;
    this.alpha = inactive;
  }
}

class Triggers extends PIXI.Graphics {
  constructor(app, labels, x, y) {
    const triggers = super();

    this.input = { [labels[0]]: 0, [labels[1]]: 0 };
    this.config = [
      { x: 0, y: 0, w: 80, h: 30, label: labels[1] },
      { x: 0, y: 50, w: 80, h: 20, label: labels[0] }
    ];

    const r = 3;

    for (const cfg of this.config) {
      const trigger = new PIXI.Graphics();

      trigger.interactive = true;
      trigger.buttonMode = true;
      trigger.alpha = inactive;

      trigger
        .lineStyle(2, line)
        .beginFill(fill)
        .drawRoundedRect(cfg.x, cfg.y, cfg.w, cfg.h, r)
        .endFill()
        .on("pointerdown", this.pointerDown.bind(this, cfg.label, trigger))
        .on("pointerup", this.pointerUp.bind(this, cfg.label, trigger))
        .on("pointerupoutside", this.pointerUp.bind(this, cfg.label, trigger));

      triggers.addChild(trigger);
    }
    app.stage.addChild(triggers);
    this.transform.position.x = x;
    this.transform.position.y = y;
  }

  pointerDown(label, trigger) {
    this.input[label] = 1;
    trigger.alpha = active;
  }

  pointerUp(label, trigger) {
    this.input[label] = 0;
    trigger.alpha = inactive;
  }
}

class Joystick extends PIXI.Graphics {
  constructor(app, x, y) {
    const outer = super();

    this.input = { x: 0, y: 0 };
    this.config = { inner: 30, outer: 80 };

    outer.interactive = true;
    outer.buttonMode = true;
    outer.alpha = inactive;

    outer
      .lineStyle(2, line)
      .beginFill(fill)
      .drawCircle(this.config.outer, this.config.outer, this.config.outer)
      .endFill()
      .on("pointerdown", this.pointerDown)
      .on("pointerup", this.pointerUp)
      .on("pointerupoutside", this.pointerUp)
      .on("pointermove", this.pointerMove);

    const inner = new PIXI.Graphics();
    inner.alpha = active;
    inner
      .lineStyle(2, line)
      .beginFill(fill)
      .drawCircle(this.config.outer, this.config.outer, this.config.inner)
      .endFill();
    this.inner = inner;

    outer.addChild(inner);
    app.stage.addChild(outer);

    this.transform.position.x = x;
    this.transform.position.y = y;
  }

  pointerDown(event) {
    this.alpha = active;
    this.start = { x: event.data.global.x, y: event.data.global.y };
    this.data = event.data;
  }

  pointerUp() {
    const shape = this.geometry.graphicsData[0].shape;
    this.alpha = inactive;

    this.data = null;
    this.input = { x: 0, y: 0 };

    this.inner.x = shape.x - shape.radius;
    this.inner.y = shape.y - shape.radius;
  }

  pointerMove() {
    if (this.data) {
      const shape = this.geometry.graphicsData[0].shape;

      const dx = this.data.global.x - this.start.x;
      const dy = this.data.global.y - this.start.y;

      const angle = Math.atan2(dy, dx);
      const distance = Math.sqrt(dx ** 2 + dy ** 2);

      const px = Math.min(distance, shape.radius) * Math.cos(angle);
      const py = Math.min(distance, shape.radius) * Math.sin(angle);

      this.inner.x = shape.x - shape.radius + px;
      this.inner.y = shape.y - shape.radius + py;

      this.input = { x: px / shape.radius, y: -py / shape.radius };
    }
  }
}

class DPad extends PIXI.Graphics {
  constructor(app, x, y) {
    const dpad = super();

    this.input = { right: 0, up: 0, down: 0, left: 0 };
    this.config = { width: 40, height: 40, radius: 5 };

    dpad.alpha = inactive;
    dpad.interactive = true;
    dpad.buttonMode = true;

    const config = this.config;

    dpad
      .lineStyle(2, line)
      .beginFill(fill, 1)
      .arc(
        config.height + config.radius,
        config.radius,
        config.radius,
        Math.PI,
        (Math.PI * 3) / 2
      )
      .arc(
        config.height + config.width - config.radius,
        config.radius,
        config.radius,
        -Math.PI / 2,
        0
      )
      .lineTo(config.height + config.width, config.height)
      .arc(
        config.height * 2 + config.width - config.radius,
        config.height + config.radius,
        config.radius,
        -Math.PI / 2,
        0
      )
      .arc(
        config.height * 2 + config.width - config.radius,
        config.height + config.width - config.radius,
        config.radius,
        0,
        Math.PI / 2
      )
      .lineTo(config.height + config.width, config.height + config.width)
      .arc(
        config.height + config.width - config.radius,
        config.height * 2 + config.width - config.radius,
        config.radius,
        0,
        Math.PI / 2
      )
      .arc(
        config.height + config.radius,
        config.height * 2 + config.width - config.radius,
        config.radius,
        Math.PI / 2,
        Math.PI
      )
      .lineTo(config.height, config.height + config.width)
      .arc(
        config.radius,
        config.height + config.width - config.radius,
        config.radius,
        Math.PI / 2,
        Math.PI
      )
      .arc(
        config.radius,
        config.height + config.radius,
        config.radius,
        Math.PI,
        (Math.PI * 3) / 2
      )
      .lineTo(config.height, config.height)
      .lineTo(config.height, config.radius)
      .endFill()
      .on("pointerdown", this.pointerDown)
      .on("pointerup", this.pointerUp)
      .on("pointerupoutside", this.pointerUp)
      .on("pointermove", this.pointerMove);

    app.stage.addChild(dpad);

    this.transform.position.x = x;
    this.transform.position.y = y;
  }

  pointerDown(event) {
    this.alpha = active;
    this.center = {
      x: this.x + this.config.height + this.config.width / 2,
      y: this.y + this.config.height + this.config.width / 2
    };
    this.data = event.data;
    this.pointerMove();
  }

  pointerUp() {
    this.alpha = inactive;
    this.data = null;
    this.input = { right: 0, up: 0, down: 0, left: 0 };
  }

  pointerMove() {
    if (this.data) {
      const dx = this.data.global.x - this.center.x;
      const dy = this.data.global.y - this.center.y;

      const angle = Math.atan2(dy, dx);
      const distance = Math.sqrt(dx ** 2 + dy ** 2);

      this.input.right = ~~(Math.abs(angle) < Math.PI / 3);
      this.input.up = ~~(Math.abs(angle + Math.PI / 2) < Math.PI / 3);
      this.input.down = ~~(Math.abs(angle - Math.PI / 2) < Math.PI / 3);
      this.input.left = ~~(
        Math.abs(angle - Math.PI) < Math.PI / 3 ||
        Math.abs(angle + Math.PI) < Math.PI / 3
      );
    }
  }
}

const init = () => {
  const app = new PIXI.Application({
    antialias: true,
    width: window.innerWidth,
    height: window.innerHeight,
    view: document.querySelector("canvas"),
    transparent: true,
    resolution: window.devicePixelRatio || 1,
    resizeTo: window,
    autoResize: true
  });

  const clientWidth = document.body.clientWidth;

  const leftJoy = new Joystick(app, 80, 30);
  const rightJoy = new Joystick(app, 0, 160);
  rightJoy.transform.position.x = clientWidth - 50 - rightJoy.config.outer * 2;

  const abxy = new ABXY(app, clientWidth - 100 - rightJoy.config.outer * 2, 10);
  const dpad = new DPad(app, 60, 220);

  const leftTriggers = new Triggers(app, ["l", "zl"], 10, 0);
  const rightTriggers = new Triggers(app, ["r", "zr"], 0, 0);
  rightTriggers.transform.position.x =
    clientWidth - 10 - rightTriggers.config[0].w;

  const plus = new Button(app, "plus", clientWidth - 300, 50);
  const minus = new Button(app, "minus", 300, 50);
  const capture = new Button(app, "capture", 300, 150);
  const home = new Button(app, "home", clientWidth - 300, 150);

  const leftJoyPress = new Button(app, "press", 300, 250);
  const rightJoyPress = new Button(app, "press", clientWidth - 300, 250);

  const socket = new WebSocket(`ws://${location.host}/controller`);

  let requestId;

  const reportInput = () => {
    const inputs = {
      dpad: dpad.input,
      button: {
        ...abxy.input,
        ...rightTriggers.input,
        ...leftTriggers.input,
        ...home.input,
        ...plus.input,
        ...minus.input,
        ...capture.input
      },
      stick: {
        left: {...leftJoy.input, ...leftJoyPress.input},
        right: {...rightJoy.input, ...rightJoyPress.input}
      }
    };
    socket.send(JSON.stringify(inputs));
    requestId = window.requestAnimationFrame(reportInput);
  };

  socket.addEventListener("open", () => {
    requestId = window.requestAnimationFrame(reportInput);
  });

  socket.addEventListener("close", () => {
    window.cancelAnimationFrame(requestId);
  });
};

const setup = () => {
  setTimeout(() => {
    if (window.innerWidth / window.innerHeight > 1) {
      init();
      window.removeEventListener("orientationchange", setup);
    }
    // Waiting for innerWidth/innerHeight change
  }, 500);
};

document.addEventListener("DOMContentLoaded", setup);
window.addEventListener("orientationchange", setup);
