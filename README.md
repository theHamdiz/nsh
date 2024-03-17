# 🚀 `nsh` previously known as **nameShift** 🚀

### 🎉 Your ultimate, quirky little buddy for **massive string transformations**! 🎉

Ever felt the urge to rename everything in sight? 😈 Whether it's those pesky filenames, the directories hiding in the shadows, or the sneaky strings lounging inside your files, **nsh** is your go-to partner in crime! 🕵️‍♀️💻

Here's what this bad boy can do for you:

- 📁 Wave goodbye to old file & directory names and say hello to their snazzy new identities.
- 📄 Dive deep into your files, leaving no string unturned.
- 🔍 Got a type? We can focus on specific file extensions for our makeover session.
- 🚀 Want speed? Go concurrent. Prefer a leisurely pace? We do synchronous, too.
- 🎯 Case-sensitive or case-agnostic, we cater to all tastes.
- 🛑 Config directories are usually off-limits, but if you're feeling rebellious, we won't stop you.
- 📊 Love numbers? Get juicy details with our tabular reports on the makeovers and mishaps.
- 🚩 Flexible with your flags, whether you like them single or double-dashed.

> No matter the challenge, **nsh** is your all-in-one, Swiss Army knife for string manipulation. Whether it's filenames, dirnames, or file guts, we've got you covered. 🛡️✨

# 🚧 **Build Instructions** 🚧

Getting **nameShift** `nsh` ready to rock is like brewing your favorite coffee - it requires specific steps but the end result is oh-so satisfying. Follow these brewing... I mean, building instructions to get started on both Windows and Unix systems. 🛠️☕

### For Windows Warriors 🛠️💻

1. Open your command prompt with a flick of your wrist or a click of your mouse. Feel the power!
2. Navigate to the root directory of **nsh** where the mystical `build.bat` script resides.
3. Summon the build script with the ancient incantation:

```bash
✅ .\build.bat
```

Watch as the console comes alive, compiling and concocting your tool into existence. It's like magic, but cooler because it's code. 🧙‍♂️💾

### For Unix Heroes 🛠️🐧

1. Open your terminal. You can pretend it's a portal to another dimension if you like - we won't judge.
2. Trek through the filesystem to the sacred land of **nsh**'s base camp, where `build.sh` has been waiting for centuries (or since you cloned the repo).
3. With the courage of a thousand penguins, execute the revered script:

```zsh
✅ ./build.sh
```

Behold as the script weaves its spells, stitching together binaries and libraries with the elegance of a master craftsman. Or craftswoman. Or craftsperson. You get the idea. 🧵✨

**🍾 Congratulations!** You've now built **nsh**, your very own digital Swiss Army knife, ready to slice and dice strings with the finesse of a gourmet chef in the digital kitchen. Go forth and refactor with reckless abandon, my friend! 🎊🔪

# 🚀 **Installing nsh, System-Wide** 🌐📦

Elevating **nameShift** `nsh` from a mere tool in your digital toolbox to a cornerstone of your system's utility belt is akin to granting it the key to the city. This step ensures `nsh` is not just another tool, but a trusted companion ready at your beck and call, across the vast landscapes of your operating system. 🗝️💼

### Universal Spell for All Environments 🌎✨

1. Arm yourself with the resolve of a seasoned administrator. You're about to grant `nsh` the honor of serving you system-wide.
2. Stand at the heart of **nsh**'s domain, where `install.py` resides, a script penned by the sages to transcend OS boundaries.
3. Invoke the installation script with the authority it requires:

> 📦 Unix Installation  

**⚠️ Note of Power:** Depending on the defensive enchantments on your system (read: permissions), you might need to invoke this script with elevated privileges. Should you encounter any resistance (errors or access denials), it's time to wield your powers as an administrator. On Unix-like systems, prepend `sudo` to the command, and on Windows, ensure your command prompt wields the might of an administrator. This is not merely a suggestion but a rite of passage for `nsh` to serve you without hindrance. 🛡️🔑

```zsh
✅ sudo python3 build/install.py
```

> Or if you have python installed at /usr/bin/env python3 just run the following directly:  

```zsh
✅ sudo ./build/install.py
```

> 📦 Windows Installation  

```bash
✅ python build\\install.py
```

> Or just run the following directly:  

```bash
✅ .\\build\\install
```

As the script springs to life, it meticulously examines the lay of the land (your OS) and lays down the pathways for `nsh` to roam freely. It doesn't just install; it creates a symbolic link (`sem-link`), a beacon for all your command-line endeavors, guiding you to `nsh` whenever the need arises. 🛣️🔗



**🎉 Et Voilà!** You've successfully bestowed upon **nsh** the honor of serving you at a system-wide level. Gone are the days of navigating to its abode for every task. `nsh` is now your ever-present aide, ready to tackle string manipulation challenges at a moment's notice. Journey forth with confidence, knowing that `nsh` is always within reach, a true testament to your command over your digital realm. 🌟🛠️

#### 🖥️ Running on Windows

Feel the power coursing through your veins as you command:

```bash
✅ .\nsh.exe "path\\to\\directory" "OldText" "NewText" --ignore-config-dirs=true -work-globally=false --concurrent-run=false -case-matching=true -file-extensions=".go,.md"
```

Or if you have installed the tool simply run it with:

```bash
✅ nsh "path\\to\\directory" "OldText" "NewText" -i=true -g=false --cr=false -cm=true --exts=".go,.md"
```

#### 🐧 On Unix Systems

Unleash the beast with:

```zsh
✅ ./nsh "path/to/directory" "OldText" "NewText" --ignore-config-dirs=true --work-globally=false -concurrent-run=false -case-matching=true --file-extensions=".go,.md"
```


Or if you have installed the tool simply run it with:

```zsh
✅ nsh "path\\to\\directory" "OldText" "NewText" --i=true -g=false --cr=false -cm=true -ext=".go,.md"
```

# 🤹 **Flexibility & Forgiveness: `nsh`'s Approach to Parameters** 🤹

In the vast and often unforgiving landscape of command-line tools, **nameShift** `nsh` stands out as a beacon of flexibility and understanding. Recognizing that the road to perfection is paved with typos and preferences, `nsh` is designed to accommodate your unique journey. Whether you favor the scenic route with `--long-form-parameters` or the express lane with `-s` shortcuts, `nsh` adapts to your style, not the other way around. 🛣️🚀

### 🌈 **A Spectrum of Choices** 🌈

- **Dual Parameter Personalities**: Each command in `nsh` has been bestowed with two personas - a verbose, descriptive one for clarity (`--ignore-config-dirs`), and a concise, shorthand alias for efficiency (`-i`). This duality ensures that whether you're a detail-oriented wizard or a speed-seeking knight, `nsh` speaks your language. 🗣️✨

- **Typos? No Problem**: Ever mistyped a parameter and faced the cold, unyielding error message of less forgiving tools? `nsh` chuckles at such rigidity. Designed with empathy, it understands the human element, accepting both `ext` and `exts` for file extensions. This gesture of understanding underscores `nsh`'s commitment to being not just a tool, but a companion on your digital odyssey. 🛠️💖

### 🎩 **Magic Behind the Curtain** 🎩

Peek behind the curtain, and you'll find `nsh`'s secret sauce - a custom flag parsing mechanism that breathes life into these user-friendly features. This mechanism is the unsung hero, allowing `nsh` to gracefully handle variations in parameter inputs without breaking a sweat. It's not just code; it's a philosophy woven into the very fabric of `nsh` - to be as adaptable and accommodating as the diverse community it serves. 🌟👥

> ✅ **nsh** is not merely a tool; it's a testament to the belief that technology should adapt to humans, not the other way around. With its blend of flexibility, forgiveness, and a touch of whimsy, `nsh` invites you to engage in the art of string manipulation with a sense of ease and joy unparallel in the realm of command-line tools. 🎭🛠️

This commitment to user experience makes `nsh` not just a powerful tool in your arsenal but a joy to use, proving that in the digital world, understanding and adaptability are the true hallmarks of innovation. 🚀💡

# 🚀 **nsh vs. The World** 🚀

In the realm of **massive string transformations**, where tools like `vidir` and `sed` have long reigned, a new challenger approaches: **nsh** `nsh`. With its quirky charm and powerful capabilities, it's not just another tool in the shed; it's your go-to gadget for turning the mundane into the extraordinary. 🌟🔧

Here's why **nsh** stands out:

- **Speed & Precision**: Unlike `vidir` and `sed`, which excel in their specific domains, `nsh` combines the best of both worlds with its speedy concurrent processing and pinpoint accuracy in string transformations. 🚀🎯
- **User-Friendly**: While `sed` commands can resemble arcane incantations, and `vidir` requires vim knowledge, `nsh` is designed with a simple and intuitive interface. Transform strings without needing to consult ancient tomes or master text editors. 📖✨
- **Versatile Functionality**: Beyond mere file renaming, `nsh` dives deep into file contents across a multitude of file extensions, making it the Swiss Army knife for all your string manipulation needs. 🛠️📁
- **Innovative Features**: With capabilities like synchronous and concurrent transformations, case-sensitive or case-agnostic operations, and detailed tabular reports, `nsh` offers a toolkit designed for modern needs. 🌐💡

### 📜 **The Road Ahead: To-Do List** 📜

As mighty as **nsh** stands, the quest for perfection is never-ending. Here are the realms yet to be conquered:

- [ ] **GUI Integration**: Bringing the power of `nsh` to a graphical user interface for those who prefer visuals over command-line spells.
- [ ] **Cross-Platform Package Managers**: Aim to distribute `nsh` through package managers like Homebrew, apt, and others, making installation a breeze across any land.
- [ ] **Advanced Pattern Matching**: Implement regex support for the adventurers who need to capture or transform more complex string patterns.
- [ ] **Localization Support**: Making `nsh` a true citizen of the world by supporting multiple languages in its interface.
- [ ] **Plugin Ecosystem**: Enabling the community to extend `nsh` with their own spells and enchantments through plugins.

> **nsh** is not merely a tool; it's a companion on your journey through the digital realm. Whether against `vidir`, `sed`, or any other, `nsh` stands tall, ready to tackle any challenge with you. 🛡️🚀

Integrate these modular additions into your existing documentation to cast a wider net, capturing the hearts of those loyal to `vidir` and `sed`, and those yet to pledge their allegiance. With **nsh**, you're not just choosing a tool; you're embracing a new ally in your endless quest for digital mastery. 🎉👑

**🍻 Here's to using it and loving it just as much as I did coding it - with a bit of sass, a dash of class, and loads of brass. Enjoy, you magnificent beast!** 🎉

