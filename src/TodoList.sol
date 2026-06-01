pragma solidity ^0.8.20;

interface ITodoToken {
    function mint(address to, uint256 amount) external;
}

contract TodoList {
    //структура 1 тудушки
    struct Todo {
        uint256 id;
        string text;
        bool isCompleted;
        uint256 completionCount; 
        bool isDeleted;
        address owner;
    }
    //инкремент, глобальный для всего контракта
    uint256 private nextId = 1;
    // маппинг map[u256]Todo
    mapping(uint256 => Todo) private todos;


    ITodoToken public token;
    uint256 public reward = 10 * 10 ** 18; // 10 токенов
    // ивенты, аналог логов
    event TodoCreated(uint256 indexed id, string text, address indexed owner);
    event TodoToggled(uint256 indexed id, bool isCompleted);
    event TodoDeleted(uint256 indexed id);
    event TodoUpdated(uint256 indexed id);
    event RewardPaid(address indexed user, uint256 amount);

    constructor(address tokenAddr) {
        token = ITodoToken(tokenAddr);
    }

    function create(string memory _text) public { // мемори - память вызова
        uint256 curId = nextId;

        todos[curId] = Todo({
            id: curId,
            text: _text,
            isCompleted: false,
            completionCount: 0,
            isDeleted: false,
            owner: msg.sender
        });

        emit TodoCreated(curId, _text, msg.sender); // вызрвать ивент

        unchecked {
            nextId++; // анчекед экономит газ
        }
    }
    function toggleComplete(uint256 _id) public {
        Todo storage todo = todos[_id]; // указатель на данные в блокчейне
        // аналог if err != nil в го, проверяем что тудушка существует и не удалена
        require(todo.id != 0, "Todo not found");
        require(!todo.isDeleted, "Todo is deleted");
        require(todo.owner==msg.sender, "Not ur todo");

        todo.isCompleted = !todo.isCompleted;
        todo.completionCount++; 
        emit TodoToggled(_id, todo.isCompleted);

        if (todo.isCompleted && todo.completionCount == 1) { 
            token.mint(msg.sender, reward);
            emit RewardPaid(msg.sender, reward);
        }
    }

    function deleteTodo(uint256 _id) public {
        Todo storage todo = todos[_id];
        require(todo.id != 0, "Todo not found");
        require(!todo.isDeleted, "Already deleted");
        require(todo.owner==msg.sender, "Not ur todo");

        todo.isDeleted = true;
        emit TodoDeleted(_id);
    }

    function getTodo(uint256 _id) public view returns (Todo memory) {
        Todo storage todo = todos[_id];
        require(todo.id!=0, "Todo not found");
        return todo;
    }
    function getAllTodos() public view returns (Todo[] memory) {
        uint256 count = 0;
        for (uint256 i = 1; i < nextId; i ++) {
            if (todos[i].owner == msg.sender && !todos[i].isDeleted) {
                count++;
            }
        }
        Todo[] memory result = new Todo[](count);
        uint256 idx = 0;
        for (uint256 i = 1; i < nextId; i++) {
            if (todos[i].owner == msg.sender && !todos[i].isDeleted) {
                result[idx] = todos[i];
                idx++;
            }
        }
        return result;
    }
    function updateText(uint256 _id, string memory _newText) public {
        Todo storage todo = todos[_id];
        require(todo.id!=0, "Todo not found");
        require(todo.owner==msg.sender, "Not ur todo");
        require(!todo.isDeleted, "Todo already deleted");

        todo.text = _newText;
        emit TodoUpdated(_id);
    }
}